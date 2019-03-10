# Catena del freddo

## Architettura
I sensori sono sviluppati in Scala e comunicano, ad intervalli regolari, i dati captati o eventuali anomalie ad un server tramite chiamata POST.
Il server, scritto in GO, conserva i dati all'interno di un DB Mongo, con una collezione per ogni  [key] DeviceID e oggetti BSON contententi [\_id; temperatura; timestamp; warning].
Tali valori sono plottati, utilizzando R, in una serie di immagini aggiornate costantemente, una per sensore. Tale immagine è visualizzabile attraverso una chiamata GET alla pagina dedicata del server, facendo riferimento al sensore corrispondente.
Nel caso in cui uno dei sensori invii una temperatura fuori norma, il server, nella pagina di monitoring del sensore, segnala l'informazione indicando il numero totale di warning ricevuti.
L'home page contiene la lista di tutti i dispositivi attivi, linkando ognuno alla rispettiva pagina, contentente il grafico della temperatura e il numero di warning ricevuti.

Il sensore è stato reso più smart. Al fine di limitare la comunicazione e l'energia spesa, il sensore non comunica le variazioni di temperatura troppo piccole in un intervallo di tempo limitato, warning esclusi.
Nel caso in cui il sensore non riesca a comunicare con il server, esso lavora in modalità offline mantenendo traccia dei warning eventualmente letti.
Periodicamente effettua dei tentativi per ristabilire la connessione e, nel momento in cui ritorna online, comunica i dati registrati al server.


### **Librerie esterne**
**GO**
https://github.com/mongodb/mongo-go-driver
```
Installare la libreria di GO inserendo nella cartella ~\go\src il contenuto di LibrerieGO.zip
```


**R**
https://jeroen.github.io/mongolite/
ggplot2
ggfortify
ggpmisc
Cairo
```
Installare le librerie tramite comando:
install.packages("mongolite")
install.packages("ggplot2")
install.packages("ggfortify")
install.packages("ggpmisc")
install.packages("Cairo")
```



## Usage

```
MongoDB deve essere installato ed essere contattabile all'indirizzo: mongodb://127.0.0.1:27017/"

Server:

Nella cartella ~\webserver\mongo eseguire il comando: go build
Nella cartella ~\server\ eseguire il comando: go run .\new-server.go
Aprire la homepage del server all'indirizzo: http://localhost:8080/

Sensore:

Buildare il sensore tramite IntelliJ IDEA: Build -> Build Artifacts -> Build
Nella cartella ~\SensoreTemperatura\out\artifacts\SensoreTemperatura_jar aprire il terminale ed eseguire il sensore tramite il comando java -jar .\SensoreTemperatura.jar
Questa chiamata accetta un parametro [sensoreID] con il quale si assegna un id al sensore. Se il parametro non è presente viene inserito l'id "default"
Eventuali parametri di configurazione possono essere modificati tramite il file config.properties presente nella stessa cartella.

Devono essere adattati alcuni valori di path.
In R-Handler.R : Path di salvataggio del jpg 	linea 8	
In new-server.go:
 PATH_RSCRIPT	path di Rscript.exe 			linea 21
 PATH_RHANDLER 	path dello script R-Handler.R 	linea 22

```
