
#!/usr/bin/python
import sys, os, subprocess as sub

class cd:
  ''' Context Manager for changing directory within python. '''

  def __init__(self, newPath):
    self.newPath = os.path.expanduser(newPath)

  def __enter__(self):
    self.savedPath = os.getcwd()
    os.chdir(self.newPath)

  def __exit__(self, etype, value, traceback):
    os.chdir(self.savedPath)


##################################################################
def answerBack(cmdin):
  rez = sub.Popen(cmdin, shell=True, stdout=sub.PIPE)
  bodz = rez.communicate()
  print "Error = ", bodz[1],
  return bodz[0]

##################################################################
if __name__ == "__main__":
  install = ["himitsu","methods","newtxes","structures"]
  buildtools = ["allcoins.go", "newbalances.go", "cretxos.go", "doublespend.go", "equilibrium.go", "history.go", "listall.go", "listone.go", "listshort.go", "mycoins.go", "sum.go", "unblock.go", "utxos.go"]
  buildmain = ["coin.go", "coinbase.go", "genesis.go", "masterpay.go", "payments.go", "ubi.go"]
  installcmd = "go install"
  buildcmd ="go build %s"

  for dir in install:
    print "cd %s and go install" % (dir),
    with cd(dir):
      bodytext = answerBack([installcmd])
    print bodytext

  for gocode in buildtools:
    builder = buildcmd % (gocode)
    print "Executing:", builder,
    bodtxt = answerBack([builder])
    print bodtxt

  for gocode in buildmain:
    builder = buildcmd % (gocode)
    print "Executing:", builder,
    bodtxt = answerBack([builder])
    print bodtxt
##################################################################
