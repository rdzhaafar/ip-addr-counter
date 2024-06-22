#!/usr/bin/env python3
import random
import math
import argparse


def gen_random_ipv4() -> str:
    ipv4 = ''
    for i in range(4):
        b = random.randint(0, 255)
        ipv4 += str(b)
        if i != 3:
            ipv4 += '.'
    return ipv4


def generate_file(max_bytes: int, filename: str):
    '''
    Generates a file with random IPv4 addresses with
    size less than or equal to max_bytes.
    '''
    # A line as at most 16 bytes
    # 255.255.255.255\n -> 15 bytes IPv4 + 1 byte for '\n'
    lines = math.floor(max_bytes / 16)
    with open(filename, 'wt') as f:
        for _ in range(lines):
            ipv4 = gen_random_ipv4()
            f.write(ipv4 + '\n')


def main():
    parser = argparse.ArgumentParser('generate.py', 'generate a file with random IPv4 addresses')
    parser.add_argument('-f', '--file', nargs='?', required=True, help='name of the generated file')
    size_group = parser.add_mutually_exclusive_group(required=True)
    size_group.add_argument('-gb', '--max-size-gb', nargs='?', type=int, help='maximum size of the generated file in GB')
    size_group.add_argument('-kb', '--max-size-kb', nargs='?', type=int, help='maximum size of the generated file in KB')
    size_group.add_argument('-mb', '--max-size-mb', nargs='?', type=int, help='maximum size of the generated file in MB')
    args = parser.parse_args()

    filename = args.file

    KB = 1024
    MB = KB * 1024
    GB = MB * 1024
    if args.max_size_kb:
        max_bytes = args.max_size_kb * KB
    elif args.max_size_mb:
        max_bytes = args.max_size_mb * MB
    elif args.max_size_gb:
        max_bytes = args.max_size_gb * GB
    else:
        raise Exception('Max size not specified')
    if max_bytes <= 0:
        raise Exception('max size cannot be negative')

    generate_file(max_bytes, filename)


if __name__ == '__main__':
    main()
