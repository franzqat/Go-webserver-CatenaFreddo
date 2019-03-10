
package Sensore

import java.net.SocketTimeoutException

import scalaj.http._

/**
  * Client per gestire le richieste HTTP del sensore: Post
  */
object HttpHandler {

  /**
    * Post al server
    * @param deviceId
    * @param timestamp
    * @param temperatura
    * @return
    */
  def post(URL: String, deviceId: String, timestamp: Long, temperatura: Double, warning: Int): Boolean = {
    try {
      Http(URL + deviceId).postForm(Seq("Device Id" -> deviceId, "timestamp" -> timestamp.toString, "temperatura" -> temperatura.toString, "warning" -> warning.toString )).asString
      return true
    }
    catch  {
      case e: SocketTimeoutException => println("Connection timed out")
        return false
      case e: Exception => println("Errore " + e)
        return false
    }

  }

}
