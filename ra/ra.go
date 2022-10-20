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

	"github.com/DistributedClocks/GoVector/govec"
)

const MAX_PROCESSES = 4

type Request struct{
    Clock   []int
    Pid     int
    Actor   cmd.ACTOR
    log []byte
}

type Reply struct{
    log []byte
}

type RASharedDB struct {
    // Constantes
    me          int     // This node's unique number
    N           int     // Number of nodes in the network
    Actor       cmd.ACTOR  // Actor type
    // Enteros
    OurSeqNum   []int     // The sequence number chosen by a request originating at this node 
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


func New(me int, usersFile string, actor_t cmd.ACTOR, logger *govec.GoLog) (*RASharedDB) {
    messageTypes := []ms.Message{Request{}, Reply{}}
    msgs := ms.New(me, usersFile, messageTypes)
    ra := RASharedDB{me, MAX_PROCESSES, actor_t, make([]int, MAX_PROCESSES), MAX_PROCESSES, false, make([]bool, MAX_PROCESSES), sync.Mutex{}, &msgs,  make(chan bool),  make(chan bool)}

    go func ()  {
        for {
            select {
            case <- ra.done:
                return
            default:
                switch msg := (ra.ms.Receive()).(type) {
                // Alguien quiere entrar en SC
                case Request:
                    logger.UnpackReceive("Receive request", msg.log, nil, govec.GetDefaultLogOptions())
                    cmd.MaxArray(ra.OurSeqNum, msg.Clock)
                    // Si no queremos entrar en SC: enviamos reply
                    // Si queremos entrar en SC enviamos reply si:
                    //      - El que quiere entrar tiene un clock mayor
                    //      - El que quiere entrar tiene un clock igual y un pid mayor
                    ra.Mutex.Lock()
                    condition := !ra.ReqCS ||
                        (cmd.Max(ra.OurSeqNum) > ra.OurSeqNum[msg.Pid-1] && cmd.Exclude(ra.Actor, msg.Actor)) ||
                        (cmd.Max(ra.OurSeqNum) > ra.OurSeqNum[msg.Pid-1] && ra.me > msg.Pid && cmd.Exclude(ra.Actor, msg.Actor))
                    ra.Mutex.Unlock()
                   
                    if condition {
                        log := logger.PrepareSend("Sending reply", nil, govec.GetDefaultLogOptions())
                        ra.ms.Send(msg.Pid, Reply{log})
                        continue
                    }

                    // Estamos en sección crítica y el mensaje es de un proceso con un clock mayor (tenemos prioridad)
                    ra.RepDefd[msg.Pid-1] = true
                    
                // Recibo respuesta/permiso para entrar en SC
                case Reply:
                    logger.UnpackReceive("Receive reply", msg.log, nil, govec.GetDefaultLogOptions())
                    if ra.ReqCS {
                        ra.OutRepCnt = ra.OutRepCnt - 1 // Permiso recibido
                        if ra.OutRepCnt == 0 {          // Todos los permisos recibidos
                            ra.chrep <- true
                        }
                    }

                // Mensaje de error
                default: return
                }
            }
        }
    }()
    
    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(logger *govec.GoLog) {
    ra.Mutex.Lock()
    ra.OurSeqNum[ra.me-1] = ra.OurSeqNum[ra.me-1] + 1
    ra.ReqCS = true
    ra.OutRepCnt = ra.OutRepCnt - 1
    ra.Mutex.Unlock()

    for i := 1; i <= ra.N; i++ {
        if i != ra.me {
            log := logger.PrepareSend("Sending request", nil, govec.GetDefaultLogOptions())
            ra.ms.Send(i, Request{ra.OurSeqNum, ra.me, ra.Actor, log})
        }
    }

    <- ra.chrep
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(logger *govec.GoLog) {
    ra.Mutex.Lock()
    ra.ReqCS = false
    ra.Mutex.Unlock()

    for j := 1; j <= ra.N; j++ {
        if ra.RepDefd[j-1] {
            ra.RepDefd[j-1] = false
            log := logger.PrepareSend("Sending reply", j, govec.GetDefaultLogOptions())
            ra.ms.Send(j, Reply{log})
        }
    }
    ra.OutRepCnt = MAX_PROCESSES
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
