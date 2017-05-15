from os import listdir
import sys

if __name__ == "__main__":

  incom = False

  if len(sys.argv) == 1:
    gofiles = [f for f in listdir(".") if f.endswith(".go")]
    for file in sorted(gofiles):
      print file
  else:
    print sys.argv[1]
    for line in open(sys.argv[1]):
      if incom:
        if line.startswith("*/"):
          sys.exit()
        else:
          print line.strip()

      if line.find("DOCSTRING") > 0:
        incom = True


