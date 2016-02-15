
set terminal png size 1024,768
set output 'gatewayandclient.png'
set autoscale
set title "gatewayandclient"
set grid
set xlabel 'Time stamp'
set ylabel 'Request'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "gatewayandclient.txt" using 1:2 title "Total Request", '' using 1:3 title "gateway1", '' using 1:4 title "gateway2", '' using 1:5 title "gateway3", '' using 1:6 title "gateway4"
replot



set output 'db.png'
set autoscale
set title "redis db"
set grid
set xlabel 'Time stamp'
set ylabel 'redis db'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "db.txt" using 1:2 title "Data conflict", '' using 1:3 title "load data", '' using 1:4 title "write data"
replot


set output 'gateway.png'
set autoscale
set title "gateway status"
set grid
set xlabel 'Time stamp'
set ylabel 'gateway'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "gateway.txt" using 1:2 title "Client Request", '' using 1:3 title "load User", '' using 1:4 title "ChatMessage",  '' using 1:5 title "http 200", '' using 1:6 title "http 503"
replot


set output 'poster.png'
set autoscale
set title "poster status"
set grid
set xlabel 'Time stamp'
set ylabel 'poster'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "poster.txt" using 1:2 title "Total PostRequest", '' using 1:3 title "PostFowardPost", '' using 1:4 title "PostFowardLocal",  '' using 1:5 title "DropedFromQueue"
replot



set output 'localposter.png'
set autoscale
set title "localposter status"
set grid
set xlabel 'Time stamp'
set ylabel 'localposter'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "localposter.txt" using 1:2 title "Total PostRequest"
#%, '' using 1:3 title "LocalFowardPost",
replot





set output 'chatmessage.png'
set autoscale
set title "chatmessage status"
set grid
set xlabel 'Time stamp'
set ylabel 'chatmessage'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "chatmessage.txt" using 1:2 title "CLientSendRequest", '' using 1:3 title "LocalPosterRecvRequest", '' using 1:4 title "LocalPosterRecvAck"
replot


