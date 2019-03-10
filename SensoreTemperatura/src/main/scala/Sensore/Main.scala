package Sensore


import scala.collection.mutable.ArrayBuffer


object Main extends App {

  var deviceId: String = _
  Configuration.CheckConfiguration() //L'oggetto Configuration si occupa della gestione delle informazione e dei parametri di configurazione del sensore
  if (args.length == 0) {
    println("No argument found, using default instead")
    deviceId = Configuration.DeviceId
  }
  else {
    deviceId = args(0)
  }
  val sleepTime = Configuration.sleepTime
  val URL = Configuration.URL

  var Sensibilita = Configuration.Sensibilita
  var TempCritica = Configuration.TempCritica

  ///Sensore SMART
  var contatoreCiclo: Int = 0
  var comparazioneThreshold: Double = 0.0
  var warning: Int = _

  ///Sensore OFFLINE
  var OfflineRecordArray = new ArrayBuffer[(Long, Double)]
  var iteratoreOfflineRecordArray: Int = 0
  val SIZE_OFFLINE_ARRAY = 50
  val NUM_MAX_CONTATORE_CICLI = 100


  var temperatura = 0.0
  var timestamp = System.currentTimeMillis / 1000 //basta la dichiarazione del tipo

  var state = ConnectionState.Online




  while (true) {
    timestamp = System.currentTimeMillis / 1000

    // L'oggetto Misura si occupa della generazione simulata della lettura del sensore di temperatura
    temperatura = Misura.LetturaTemperatura()

    //il sensore è una macchina a stati, si comporta in maniera diversa se la connessione è online o offline
    state match {
      case ConnectionState.Online =>
        //confronta la temperatura con la temperatura critica e restituisce 0 se il valore è sotto la soglia, 1 se è sopra
        warning = checkWarning()


        //invia la comunicazione ogni 5 controlli a meno di una variazione significativa
        //la presenza di un warning viene segnalata in ogni caso
        if ((math.abs(temperatura - comparazioneThreshold) < Sensibilita) && contatoreCiclo < 5 && warning == 0) {
          contatoreCiclo += 1
          println("Variazione piccola, non eseguo la POST. La variazione corrente rispetto all'ultimo valore postato è " + BigDecimal(math.abs(temperatura - comparazioneThreshold)).setScale(2, BigDecimal.RoundingMode.HALF_UP).toDouble)
        }
        else {
          contatoreCiclo = 0 //reset contatore
          comparazioneThreshold = temperatura
          println("Sensore online: invio " + deviceId + " " + timestamp + " " + temperatura + " warn:" + warning )
          //post al server
          if (tentativoConnessione()) { //esegue un tentativo di connessione
          } else {
            println("Cambio di stato : Disconnesso, lavoro in offline")
            vaiOffline()
          }
        }
        Thread.sleep(sleepTime)


      case ConnectionState.Offline => println("Offline")
        if (contatoreCiclo < NUM_MAX_CONTATORE_CICLI) {
          contatoreCiclo += 1

          if (checkWarning() == 1) {

            //se c'è un warning aggiunge al vettore OfflineRecordArray la misura corrente
            println("Trovato un warning, inserisco nel vettore record offline all'indice: " + iteratoreOfflineRecordArray)

            if (OfflineRecordArray.length > SIZE_OFFLINE_ARRAY) { //il sensore conserva un massimo di SIZE_OFFLINE_ARRAY elementi. Se viene superato sostituisce il nuovo elemento con il meno recente
              OfflineRecordArray.remove(iteratoreOfflineRecordArray)
            }
            OfflineRecordArray.insert(iteratoreOfflineRecordArray, (timestamp, temperatura)) //inserisco una tupla di (timestamp, temperatura)
            iteratoreOfflineRecordArray += 1 //sposto l'iteratore
            iteratoreOfflineRecordArray %= SIZE_OFFLINE_ARRAY //sposto l'iteratore in maniera circolare
          }

        }
        else {
          //dopo aver aspettato NUM_MAX_CONTATORE_CICLI prova ad eseguire il tentativo di connessione
          if (tentativoConnessione()) {
            println("Riconnessione eseguita")

            //Se non è vuoto, invia il contenuto di OfflineRecordArray
            if (OfflineRecordArray.nonEmpty) {
              println("Invio dei valori di warning registrati in modalita' offline")
              for ((timestamp, temperatura) <- OfflineRecordArray) {
                HttpHandler.post(URL, deviceId, timestamp, temperatura, 1)
                Thread.sleep(100) //per evitare il blocco
              }
              OfflineRecordArray.clear() //dopo aver inviato il contenuto, il vettore viene svuotato
            }
            vaiOnline() //cambio di stato

          } else {
            println("Connessione fallita, tentativo in " + sleepTime * 100 + " secondi")
            contatoreCiclo = 0
          }
        }
        Thread.sleep(sleepTime)
    }


  }


  /**
    * esegue fino a 5 tentativi di connessione
    * restituisce true se la connessione ha successo, false in caso contrario
    * @return
    */
  def tentativoConnessione(): Boolean = {

    for( i <- 0 to 4) {
      if (HttpHandler.post(URL, deviceId, timestamp, temperatura, warning)) {
        return true
      }
      else
        println("Retrying in " + sleepTime + " seconds. . .")
        Thread.sleep(sleepTime)
    }
    if (HttpHandler.post(URL, deviceId, timestamp, temperatura, warning)) {
      return true
    }
    return false
  }

  /**
    * Ritorna 1 se trova un warning
    * @return
    */
  def checkWarning(): Int = {
    if (temperatura < TempCritica) {
      return 0
    }
    else {
      return 1
    }

  }

  def vaiOnline() = {
    state = ConnectionState.Online
    contatoreCiclo = 0
    comparazioneThreshold = 0.0
    warning = 0
  }

  def vaiOffline() = {
    state = ConnectionState.Offline
    contatoreCiclo = 0
    comparazioneThreshold = 0.0
    warning = 0
  }

}
