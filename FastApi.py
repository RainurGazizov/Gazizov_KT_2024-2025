import shutil
from typing import List, Optional
from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.responses import FileResponse
from fastapi.staticfiles import StaticFiles
import os
import uuid
from PIL import Image
import ffmpeg
import mimetypes


app = FastAPI()

#Папка для хранения отправленных на сервер фото или видео файлов
app.mount('/media', StaticFiles(directory='media'), name='media')

#Отправка файла на сервер
@app.post("/api/upload/")
async def root(file: UploadFile = File (...)):
    if  file.content_type.split('/')[0] == "image" or file.content_type.split('/')[0] == "video":
        path = f'media/{file.filename}'
        with open(path, "wb+") as buffer:
            shutil.copyfileobj(file.file, buffer)
        
        id = str(uuid.uuid4())
        
        #Сохранение ID, имени файла и тип файла в текстовый документ
        with open('ID.txt', 'a+') as f:
            f.write(f"{id}*media/{file.filename}*{file.content_type.split('/')[0]}\n")
        file.close()
        
        return {"file_ID": id, 
                "file_type": file.content_type, 
                "file_size": file.size}
    else:
        raise HTTPException(status_code=404, detail='Incorrect file type!')


#Вывод всех ID
#@app.get("/api/upload/ID/")
#async def get_ID():
#    with open('ID.txt', 'r') as file:
#        file.seek(0)
#        lines = [line.split('*')[0] for line in file]
#    return lines


#@app.get('/api/file_download/{id}', response_class = FileResponse)
#async def get_file_download(id: str):
#    path = f'media/{id}'
#    return path


#Скачивание файла
@app.get('/api/file/{id}', response_class = FileResponse)
async def get_file(id: str, width: Optional[int]=None, height: Optional[int]=None):
    #path = ""
    with open('ID.txt', 'r') as file:
        file.seek(0)
        for line in file:
            if line.split('*')[0] == id:
                path = line.split('*')[1]
                type = line.split('*')[2].split('\n')[0]
                if (width == None) and (height == None):
                    return path
                elif type == "video":
                    ffmpeg.input(path).output(f'media/{id}_frame_%d.png', vframes=1).run()
                    path2 = f'media/{id}_frame_1.png'
                    img = Image.open(path2)
                    new_image = img.resize((width, height))
                    path = path2
                    new_image.save(path)
                    return path   
                else: 
                    img = Image.open(path)
                    new_image = img.resize((width, height))
                    path = f'media/{width}and{height}{path.split('/')[1]}'
                    new_image.save(path)
                    return path
        
    raise HTTPException(status_code=404, detail='File not found!')
    

#@app.post("/img")
#async def upload_image (files: List [UploadFile]= File(...)): 
#    for img in files:
#        with open(f'{img.filename}', "wb") as buffer:
#            shutil.copyfileobj(img.file, buffer)
#    return {"file_name": "Good"}