import os
import pytesseract
from PIL import Image

def remove_invalid_char(obj):
    """去除字符串中无效字符
    input: str
    out: str
    """
    vaild_char = [char for char in obj if char.isidentifier()]

    return ''.join(vaild_char)


def ocr_extract_word(img_path):
    """ocr提取图片内文字
    input: img_path str
    out: str
    """
    img = Image.open(img_path)
    out = pytesseract.image_to_string(img, lang="chi_sim")

    return remove_invalid_char(out)


if __name__ == "__main__":
    path = '/home/yodark/Pictures/'
    for dirpath, dirnames, filenames in os.walk(path):
        for filename in filenames:
            if filename.split('.')[-1] not in ['jpeg', 'png', 'jpg']:
                continue
            image = os.path.join(dirpath, filename)

            out = ocr_extract_word(image)
            print(filename, out)
