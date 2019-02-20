args = commandArgs(trailingOnly=TRUE)
if (length(args)==0) {
  stop("Inserisci come argomento l'id del sensore", call.=FALSE)
}

#path <- paste("D://magistrale/example1.jpg", sep="")
path <- paste("C://Users/franz/go/src/webserver/server/sensori/",args[1],"/",args[1],".jpg", sep="")

library(mongolite)
library(ggplot2)
library(scales) # per date_format
library(ggpmisc)
library(Cairo)
options(device="CairoWin")
dmd <- mongo(args[1], url = "mongodb+srv://utente:unict@progettoapl-zkgjt.mongodb.net/test?retryWrites=true")
#45588774
alldata <- dmd$find(sort = '{"timestamp": -1}', limit = 100)

 #print(alldata)


# converte i timestamp in orario leggibile
alldata$timestamp <- as.POSIXct(as.numeric(as.character(alldata$timestamp)), origin="1970-01-01", tz="Etc/GMT+1")

#PLOT

alertdata <- subset(alldata, as.numeric(temperatura) > -19.5)


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
if(nrow(alertdata)> 0){ 
  p+
  geom_point(data=alertdata, aes(x=timestamp, y=as.numeric(temperatura)), colour="red", size=5)+
  geom_text(angle =45, col="red", data= alertdata, vjust = -0.5, hjust = -0.4, aes(label=as.numeric(temperatura)))
}
 
  ggsave(filename = path, type="cairo", width=10, height=10, units="in", dpi=150)
  q()


