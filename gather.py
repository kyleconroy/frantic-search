import os
import requests
import multiprocessing

CARD_URL = "http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid={}"

def scrape(queue):
    while True:
        multiverse_id = queue.get()
        print(multiverse_id)
        queue.put(multiverse_id + multiprocessing.cpu_count())


if __name__ == "__main__":
    try:
        os.mkdir("cards")
    except OSError:
        pass

    manager = multiprocessing.Manager()
    queue = manager.Queue()
    pool = multiprocessing.Pool(multiprocessing.cpu_count())

    for i in range(1, 5):
        queue.put(i)

    i = pool.map_async(scrape, (queue,))
    i.wait()
