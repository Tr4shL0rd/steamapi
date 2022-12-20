import csv
import json
import requests
import os
from dotenv import load_dotenv
load_dotenv()
def remove_dupes(lst):
    unique_list = []
    for elem in lst:
        if elem not in lst:
            if elem not in unique_list:
                unique_list.append(elem)
    return unique_list

api_key = os.environ["API_KEY"]
steam_id = os.environ["STEAM_ID"]
app_id = "244850"

url = f"https://api.steampowered.com/ISteamApps/GetAppList/v2/?key={api_key}"
r = requests.get(url)
data = r.json()
header = ["appid", "name"]
app = [[dat["appid"] or "None", dat["name"] or "None"] for dat in data["applist"]["apps"]]
app = remove_dupes(app)
with open("apps.csv", "w", newline="") as f:
    writer = csv.DictWriter(f, fieldnames=header)
    writer.writeheader()
    for row in app:
        writer.writerow({"appid": row[0], "name": row[1]})