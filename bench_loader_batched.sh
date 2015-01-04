for i in 10 100 200 500 1000 2000 5000 10000 100000; do
 rm bolt.db; echo -n $i,; /usr/bin/time -f "%E" ./loader_batched -wb $i
done
