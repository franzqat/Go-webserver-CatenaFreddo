package Sensore

import scala.util.Random

object Misura {


  /**
    *Mediana della distribuzione normale
    */
  var piccoGaussiana : Double = -22


  def LetturaTemperatura() = SimulaTemperatura()


  /**
    * Viene simulato un andamento della temperatura con una variazione dovuta ad un piccolo numero casuale
    * NB la catena del freddo è rotta se la temperatura supera i -18 gradi
    * @return temperatura
    */
  def SimulaTemperatura() = {
    piccoGaussiana = piccoGaussiana + Random.nextGaussian()/10
    BigDecimal(piccoGaussiana).setScale(3, BigDecimal.RoundingMode.HALF_UP).toDouble
  }
}
