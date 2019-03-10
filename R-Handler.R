args = commandArgs(trailingOnly=TRUE)
if (length(args)==0) {
  stop("Inserisci come argomento l'id del sensore", call.=FALSE)
}

#path <- paste("D://magistrale/example1.jpg", sep="")
#path hardcoded dove salvare il jpg
path <- paste("C://Users/franz/go/src/webserver/server/sensori/",args[1],"/",args[1],".jpg", sep="")

library(mongolite) #libreria per mongo
library(ggplot2) #libreria per plot
library(scales) # per date_format
library(ggpmisc) #per plot picchi e valli

library(Cairo) #libreria grafica anti aliasing
options(device="CairoWin")

#test è il nome del db
connessioneMongo <- mongo(args[1], url = "mongodb://127.0.0.1:27017/test")
#connessioneMongo <- mongo(args[1], url = "mongodb+srv://utente:unict@progettoapl-zkgjt.mongodb.net/test?retryWrites=true")  #vecchia stringa x cloud

#ad all data sono assegnati gli ultimi 100 valori ordinati per timestamp in ordine crescente
alldata <- connessioneMongo$find(sort = '{"timestamp": -1}', limit = 100)

#print(alldata)

# La rimozione dell'elemento effettua automaticamente la disconnessione da mongo
rm(connessioneMongo)

# converte i timestamp in orario leggibile
alldata$timestamp <- as.POSIXct(as.numeric(as.character(alldata$timestamp)), origin="1970-01-01", tz="Etc/GMT+1")

#Subset di warning, -18 è la temp corretta
alertdata <- subset(alldata, as.numeric(temperatura) > -18.0)

#PLOT
p <- ggplot(alldata, aes(x=timestamp, y=as.numeric(temperatura), group=1)) +
  geom_line(size=0.8, colour = "black") +
  ylim(-26,-16) +
  stat_valleys(colour = "blue", span=50, size=2)       +
  stat_valleys(geom = "text", colour = "blue", angle = 45,
  vjust = 1.5, hjust = 1,  aes(label=as.numeric(temperatura), size=1), span=50)     +
  stat_peaks(colour = "orange",span=50, size=2)     +
  stat_peaks(geom = "text", angle =45, col="red",
  vjust = -0.5, hjust = -0.4, aes(label=as.numeric(temperatura), size=1), span=50) +
  geom_hline(yintercept=-18, colour="red", size=0.5)+
  labs(title = "Rivelazione Temperature\n", x = "Orario", y = "Temperatura", color = "Legend Title\n") +
  theme_bw() +
  theme(axis.text.x = element_text(size = 14), axis.title.x = element_text(size = 16),
  axis.text.y = element_text(size = 14), axis.title.y = element_text(size = 16),
  plot.title = element_text(size = 20, face = "bold", color = "darkgreen")) +
  theme(legend.position="topright")+
  scale_x_datetime(labels = date_format("%m-%d  %H:%M")) +
  theme(axis.text.x = element_text(angle = 0, vjust=0.5, size=10),panel.grid.minor = element_blank())
#Se entra nell'if ci sono dei valori critici che vengono evidenziati nel plot
if(nrow(alertdata)> 0){ 
  p+
  geom_point(data=alertdata, aes(x=timestamp, y=as.numeric(temperatura)), colour="red", size=5)+
  geom_text(angle =45, col="red", data= alertdata, vjust = -0.5, hjust = -0.4, aes(label=as.numeric(temperatura)))
}
  #salvataggio del plot in formato jpg
  ggsave(filename = path, type="cairo", width=10, height=10, units="in", dpi=150)

  q()


