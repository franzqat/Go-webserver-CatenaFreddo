# Go-Webserver

## **Librerie esterne**
**GO**
https://github.com/mongodb/mongo-go-driver

**R**
https://jeroen.github.io/mongolite/
ggplot2
ggfortify
ggpmisc
Cairo


## Architettura
Viene associato un sensore di temperatura ad ogni trasporto.
I sensori sono sviluppati in Scala e comunicano, ad intervalli regolari, i dati captati o eventuali anomalie ad un server tramite interfaccia REST.
Il server, scritto in GO, conserva i dati all'interno di un DB Mongo (con una collezione per ogni  [key] DeviceID e oggetti BSON contententi [\_id; temperatura; timestamp; warning]).
Tali valori sono plottati, utilizzando R, in una serie di immagini aggiornate costantemente, una per sensore. Tale immagine è visualizzabile attraverso una chiamata GET con una interfaccia REST sulla pagina dedicata del server, facendo riferimento al sensore corrispondente al trasporto.
Nel caso in cui uno dei sensori invii una temperatura fuori norma, il server, nella pagina di monitoring del sensore, segnala l'informazione indicando il numero di warning ricevuti.
L'home page contiene la lista di tutti i dispositivi attivi, linkando ognuno alla rispettiva pagina, contentente il grafico della temperatura e il numero di warning ricevuti.

Il sensore è stato reso più smart. Al fine di limitare la comunicazione e l'energia spesa, il sensore non comunica le variazioni di temperatura troppo piccole in un intervallo di tempo limitato, warning esclusi.
Nel caso in cui il sensore non riesca a comunicare con il server, esso lavora in modalità offline mantenendo traccia dei warning eventualmente letti.
Periodicamente effettua dei tentativi per ristabilire la connessione e, nel momento in cui ritorna online, comunica i dati registrati al server.

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

## Behaviour

Optimal PATH:
Premesse: 
	Il server è attivo e in ascolto
	Il database MongoDB è attivo
	L'id dei sensori è univoco
	
1. Il sensore periodicamente invia il valore di temperatura [simulato], il timestamp e il proprio ID attraverso una chiamata POST al server GO
2. Il server riceve la POST e: 
	2.1. Inserisce i dati nel database MongoDB nella collezione corrispondente all'ID
	2.2. Crea, se non esistono, le cartelle e i files associati al sensore con quell'ID
	2.3. Esegue lo script di R associato all'ID
3. R elabora i dati e genera il jpg corrispondente e lo conserva nella cartella del server associata all'ID del sensore

Dal punto di vista dell'utente:
L'utente accede alla pagina localhost:8080 dove trova un elenco di links chiamate con l'ID del sensore.
I link portano ad una pagina contente l'immagine plottata.