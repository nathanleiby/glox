
with open('tokens.txt', 'r') as f:
    for l in f.readlines():
        l = l.strip()
        if l == "":
            continue
        print("{}: \"{}\",".format(l, l))
