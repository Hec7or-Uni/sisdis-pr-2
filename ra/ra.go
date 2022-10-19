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
	"strconv"
	"sync"

	"github.com/DistributedClocks/GoVector/govec"
)

type Request struct{
    Clock   int
    Pid     int
    Actor   cmd.ACTOR
    VCM     []byte  // Vector Clock Message (GoVector)
}

type Reply struct{
    VCM     []byte  // Vector Clock Message (GoVector)
}

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
    // Logger
    logger      *govec.GoLog
}


func New(me int, usersFile string, actor_t cmd.ACTOR) (*RASharedDB) {
    messageTypes := []ms.Message{Request{}, Reply{}}
    msgs := ms.New(string(actor_t) + strconv.Itoa(me), "log_" + strconv.Itoa(me), messageTypes)
    logger := govec.InitGoVector(me, "LogFile", govec.GetDefaultConfig())
    ra := RASharedDB{me, 2, actor_t, 0, 0, 2, false, make([]bool, 2), sync.Mutex{}, &msgs,  make(chan bool),  make(chan bool), logger}

    go func ()  {
        for {
            select {
            case <- ra.done:
                return
            default:
                switch msg := (ra.ms.Receive()).(type) {
                // Alguien quiere entrar en SC
                case Request:

                    var response []byte
                    ra.logger.UnpackReceive("Solicitud de acceso a SC", msg.VCM, &response, govec.GetDefaultLogOptions())
                    // Si no queremos entrar en SC: enviamos reply
                    // Si queremos entrar en SC enviamos reply si:
                    //      - El que quiere entrar tiene un clock mayor
                    //      - El que quiere entrar tiene un clock igual y un pid mayor
                    ra.Mutex.Lock()
                    ra.HigSeqNum = cmd.Max(ra.HigSeqNum, msg.Clock)
                    condition := !ra.ReqCS ||
                        (ra.HigSeqNum > msg.Clock && cmd.Exclude(ra.Actor, msg.Actor)) ||
                        (ra.HigSeqNum == msg.Clock && ra.me > msg.Pid && cmd.Exclude(ra.Actor, msg.Actor))
                    ra.Mutex.Unlock()
                   
                    if condition {
                        payload := []byte("Permito acceso a SC")
                        VCM := ra.logger.PrepareSend("Permitiendo acceso a SC", payload, govec.GetDefaultLogOptions())
                        ra.ms.Send(msg.Pid, Reply{VCM})
                        continue
                    }

                    ra.logger.LogLocalEvent("Encolo peticiones de acceso a SC", govec.GetDefaultLogOptions())
                    // Estamos en sección crítica y el mensaje es de un proceso con un clock mayor (tenemos prioridad)
                    ra.RepDefd[msg.Pid-1] = true
                    
                // Recibo respuesta/permiso para entrar en SC
                case Reply:
                    var response []byte
                    ra.logger.UnpackReceive("Permiso de acceso a SC recibido", msg.VCM, &response, govec.GetDefaultLogOptions())

                    if ra.ReqCS {
                        ra.OutRepCnt = ra.OutRepCnt - 1 // Permiso recibido
                        if ra.OutRepCnt == 0 {          // Todos los permisos recibidos
                            ra.logger.LogLocalEvent("Entro en SC", govec.GetDefaultLogOptions())
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
func (ra *RASharedDB) PreProtocol(){
    ra.Mutex.Lock()
    ra.OurSeqNum = ra.OurSeqNum + 1
    ra.ReqCS = true
    ra.OutRepCnt = ra.OutRepCnt - 1
    ra.Mutex.Unlock()

    for i := 1; i <= ra.N; i++ {
        if i != ra.me {
            payload := []byte("Solicitud de SC")
            VCM := ra.logger.PrepareSend("Solicitando acceso a SC", payload, govec.GetDefaultLogOptions())
            ra.ms.Send(i, Request{ra.OurSeqNum, ra.me, ra.Actor, VCM})
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

    for j := 1; j <= ra.N; j++ {
        if ra.RepDefd[j-1] {
            ra.RepDefd[j-1] = false
            payload := []byte("Abandono la SC")
            VCM := ra.logger.PrepareSend("Notificando salida de SC", payload, govec.GetDefaultLogOptions())
            ra.ms.Send(j, Reply{VCM})
        }
    }
    ra.OutRepCnt = 2
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
