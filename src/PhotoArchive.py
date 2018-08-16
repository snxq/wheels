# @author: snxq
import datetime
import os
import os.path
import shutil
import sys
import time
"""
照片存档
参数：源文件目录，存档目录
"""

def existsOrCreate(path):
    if not os.path.exists(path):
        os.makedirs(path)


def archive(src_path, save_path):
    count = 0
    for dirpath, dirnames, filenames in os.walk(src_path):
        for filename in filenames:
            count += 1
            src_file = os.path.join(dirpath, filename)

            last_modified = os.path.getmtime(src_file)
            last_modified_datetime = datetime.datetime.fromtimestamp(last_modified)

            dst_filepath = os.path.join(
                save_path, last_modified_datetime.strftime('%Y%m%d')
            )
            dst_file = os.path.join(dst_filepath, filename)

            existsOrCreate(dst_filepath)
            if not os.path.exists(dst_file):
                shutil.copyfile(src_file, dst_file)
                print(f'{count}\t{filename} copy completed.')
            else:
                print(f'{count}\t{filename} already exists.')


if __name__ == "__main__":
    try:
        src_path = sys.argv[1]
        save_path = sys.argv[2]
    except:
        src_path = '~/100D3400'
        save_path = '~/Archiving'
    archive(src_path, save_path)
