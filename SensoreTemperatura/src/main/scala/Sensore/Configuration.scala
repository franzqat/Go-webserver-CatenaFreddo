package Sensore

import java.io.{File, FileOutputStream, IOException}
import java.nio.file.{Files, Paths}
import java.util.Properties

import com.typesafe.config.ConfigFactory

/**
  * Assegna i valori di configurazione leggenddoli dal file di configurazione
  */
object Configuration {

  var DeviceId : String = _
  var sleepTime : Int = _
  var URL : String = _
  var Sensibilita : Double = _
  var TempCritica : Double = _

  /**
    * Crea una configurazione di default se non esiste
    */
  def SaveDefaultConfiguration(): Unit = {

    var output : FileOutputStream =  new FileOutputStream("config.properties")
    val p: Properties = new Properties()

    try {
      // valori di default non compresi
      p.setProperty("deviceID", "default")
      p.setProperty("sleepTime", "5000")
      p.setProperty("URL", "http://localhost:8080/save/")
      p.setProperty("Sensibilita", "0.5")
      p.setProperty("TempCritica", "-18.0")

      p.store(output, null)
    } catch {
      case _: IOException => println("Errore nella scrittura del file")

    } finally if (output != null) try
      output.close
    catch {
      case _: IOException => println("Errore nella chiusura del file")

    }
  }

  /**
    * Controlla se esiste la configurazione
    */
  def CheckConfiguration() {
    if (!Files.exists(Paths.get("config.properties"))) {
      println("Creating default configuration")
      SaveDefaultConfiguration()
    }
    LoadConfiguration() //la carica se esiste
  }

    def LoadConfiguration(): Unit = {
      println("Loading configuration")
      val config = ConfigFactory.parseFile(new File("config.properties"))
      DeviceId = config.getString("deviceID")
      sleepTime = config.getString("sleepTime").toInt
      URL = config.getString("URL")
      Sensibilita = config.getString("Sensibilita").toDouble
      TempCritica = config.getString("TempCritica").toDouble
    }

}


