import pika, sys, os, json
import redis
import pickle
from pymongo import MongoClient

client = MongoClient('mongodb://35.192.86.242:27017/')  # Dbmongo
db = client.pythonmongodb

r = redis.Redis(host='35.192.86.242', port=6379) # Dbredis

def main():

    connection = pika.BlockingConnection(pika.ConnectionParameters(host='rabbitmq'))
    channel = connection.channel()

    channel.queue_declare(queue='covid-cases')

    def callback(ch, method, properties, body):
        caso = json.loads(body.decode('utf-8'))
        print(type(caso))
        # print(caso)
        try:
            result = db.casos.insert_one(caso)
            print('Caso: {} insertado'.format(result.inserted_id))

            print(str(result.inserted_id))
            # p_caso = pickle.dumps(caso)               # codifica el diccionario
            r.lpush("casos", body)                      # inserta el dato
            # value = r.get(str(result.inserted_id))    # obtengo el dato codificado
            # print(pickle.loads(value))
        except:
            print('Error al insertar dato.')

    channel.basic_consume(queue='covid-cases', on_message_callback=callback, auto_ack=True)
    print("Esperando por mensajes...")
    channel.start_consuming()


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print('Interrupted')
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)