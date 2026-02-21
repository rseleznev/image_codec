# RMZ фото-кодек

Примитивный фото-кодек, сделан исключительно в тренировочных целях

На текущий момент реализовано:
* кодирование массива байт [R, G, B, ...]
  * дельта-кодирование
  * RLE-кодирование
  * коды Хаффмана
* сериализация закодированных данных
* сохранение в файл
* чтение и проверка файла
* десериализация
* декодирование в исходные байты
* открытие файлов .rmz в кастомном просмотрщике

Ограничения:
* размер изображения не более 1000х1000 пикселей

Планы на будущее:
* Преобразование в модель YCbCr
* Применение дискретного косинусного преобразования
* Применение квантования
* Применение зигзаг-сканирования

## CLI-команды

* encode -width={} -height={} -inputFile={} -outputFile={}
  * Кодирует байты из inputFile (нужен RAW файл) в размерах width/height в outputFile
* decode -sourceFile={}
  * Декодирует байты из sourceFile (.rmz файл), восстанавливает исходные байты
* view -sourceFile={}
  * Открывает указанный .rmz файл в кастомном просмотрщике

Пример команды на открытие изображения в просмотрщике:
```bash
image_codec view -sourceFile="/path/image.rmz"
```

## Использование на Ubuntu

* Создать .desktop файл в ~/.local/share/applications/rmz-viewer.desktop

```
[Desktop Entry]
Type=Application
Name=RMZ Image Viewer
Comment=Viewer for RMZ image format
Exec=/полный/путь/к/image_codec view -sourceFile=%f
Icon=путь к иконке
Terminal=false
MimeType=image/x-rmz;
Categories=Graphics;Viewer;
StartupNotify=true
```

* Создать MIME-тип в ~/.local/share/mime/packages/rmz.xml

```xml
<?xml version="1.0" encoding="UTF-8"?>
<mime-info xmlns="http://www.freedesktop.org/standards/shared-mime-info">
  <mime-type type="image/x-rmz">
    <comment>RMZ encoded image</comment>
    <glob pattern="*.rmz"/>
    <icon name="image-x-generic"/>
  </mime-type>
</mime-info>
```

* Обновить базу MIME-типов

```bash
update-mime-database ~/.local/share/mime
```

* Установить ассоциацию файлов

```bash
xdg-mime default rmz-viewer.desktop image/x-rmz
```

* Конвертировать JPEG в RAW

```bash
convert input.jpg -depth 8 rgb:output.raw
```