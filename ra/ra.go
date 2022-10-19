/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
 */
package ra

import (
	"sisdis-pr-2/cmd"
	"sisdis-pr-2/ms"
	"sync"
)

type Request struct{
    Clock   int
    Pid     int
    Actor   cmd.ACTOR
}

type Reply struct{}

type RASharedDB struct {
    // Constantes
    me          int     // This node's unique number
    N           int     // Number of nodes in the network
    Actor       cmd.ACTOR  // Actor type
    // Enteros
    OurSeqNum   int     // The sequence number chosen by a request originating at this node 
    HigSeqNum   int     // The highest sequence number seen in any REQUEST message sent or recived
    OutRepCnt   int     // The number of REPLY  messages still expected
    // Booleanos
    ReqCS       bool    // True if the node is requesting the critical section
    RepDefd     []bool  // The reply_deferred[j] is TRUE when this node is deferring a REPLY to j's REQUEST message
    // Semaforo binario
    Mutex       sync.Mutex  // Mutex para proteger concurrencia sobre las variables
    // Otros
    ms          *ms.MessageSystem
    done        chan bool
    chrep       chan bool
}


func New(me int, usersFile string, actor_t cmd.ACTOR) (*RASharedDB) {
    messageTypes := []ms.Message{Request{}, Reply{}}
    msgs := ms.New(me, usersFile, messageTypes)
    ra := RASharedDB{me, 2, actor_t, 0, 0, 0, false, make([]bool, 2), sync.Mutex{}, &msgs,  make(chan bool),  make(chan bool)}

    go func ()  {
        for {
            select {
            case <- ra.done:
                return
            default:
                data := ra.ms.Receive()
                switch msg := data.(type) {
                // Alguien quiere entrar en SC
                case Request:
                    // Si no queremos entrar en SC: enviamos reply
                    // Si queremos entrar en SC enviamos reply si:
                    //      - El que quiere entrar tiene un clock mayor
                    //      - El que quiere entrar tiene un clock igual y un pid mayor
                    ra.Mutex.Lock()
                    condition := !ra.ReqCS ||
                        ra.HigSeqNum > msg.Clock && cmd.Exclude(ra.Actor, msg.Actor) ||
                        ra.HigSeqNum == msg.Clock && ra.me > msg.Pid && cmd.Exclude(ra.Actor, msg.Actor)
                    ra.Mutex.Unlock()
                   
                    if condition {
                        ra.ms.Send(msg.Pid, Reply{})
                        continue
                    }
                    
                    // Estamos en sección crítica y el mensaje es de un proceso con un clock mayor (tenemos prioridad)
                    ra.RepDefd[msg.Pid] = true
                    
                // Recibo respuesta/permiso para entrar en SC
                case Reply:
                    if ra.ReqCS {
                        ra.OutRepCnt = ra.OutRepCnt - 1 // Permiso recibido
                        if ra.OutRepCnt == 0 {          // Todos los permisos recibidos
                            ra.chrep <- true
                        }
                    }

                // Mensaje de error
                default: continue
                }
            }
        }
    }()
    
    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(){
    ra.Mutex.Lock()
    ra.OurSeqNum = ra.OurSeqNum + 1
    ra.ReqCS = true
    ra.OutRepCnt = ra.OutRepCnt - 1
    ra.Mutex.Unlock()

    for i := 0; i < ra.N; i++ {
        if i != ra.me {
            ra.ms.Send(i, Request{ra.OurSeqNum, ra.me, ra.Actor})
        }
    }

    <- ra.chrep
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(){
    ra.Mutex.Lock()
    ra.ReqCS = false
    ra.Mutex.Unlock()

    for j := 0; j < ra.N; j++ {
        if j != ra.me && ra.RepDefd[j] {
            ra.ms.Send(j, Reply{})
        }
    }
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
