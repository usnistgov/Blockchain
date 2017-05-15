import sys
import re
import os
import urllib2 as earl
import time

#==========================================================================
#This is a file of helper functions.  No main routine.
#def printerror(stringin):
#def urlpeel(ceorl):
#def printlist(listo, nested=0):
#def printfancy(listo):
#def printdict(dicto):
#def tabify(listo):
#def stringify(listo):
#def cherrypick(listin, mode):
#def addunique(target, source):
#def inlist(tolist, elem):
#def findall(onestring, alist):
#def savefile(alist, filesaver, mode):
#def savedict(adict, filesaver, mode):
#def filegetter(filename):
#def stripper(inlist):
#def editstring(atup, after):
#def visit(arg, dirname, names):
#def getconfig(configfile):
#==========================================================================
#Print to error output:
def printerror(stringin):
  sys.stderr.write("%s\n" % (stringin))

#open the url and return its contents as a list of lines.
def urlpeel(ceorl):
  lines = []
  try:
    parfp = earl.urlopen(ceorl)
    lines = parfp.readlines()
  except:
    sys.stderr.write("urlpeel: failed to open %s\n" % (ceorl))
  return lines

#supporting functions: #printlist
#print the contents of a list:
def printlist(listo, nested=0):
  if len(listo) == 0:
    print "empty list"
    return

  for url in listo:
    if url.strip() == "":
      continue

    if nested == 0:
      print "%s" % (url.strip())
    else:
      print tabify(url)


#print the contents of a 2D list:
def print2dlist(listo):
  if len(listo) == 0:
    print "empty list"
    return

  for url in listo:
    for el in url: print el

#Get unigrams from a 2D list:
def unigrams(twolist):
  for url in listo:
    for el in url: 
      if isalpha(el[0]):	#string test
        wds = el.split()


#Convert the lsit of strings 'listo' into a list of lists 'listed'.
def listify(listo):
  listed = []
  if len(listo) == 0:
    print "empty list"
    return []

  for urlo in listo:
    url = urlo.strip()
    if url == "":
      continue
    if url.find("[") == 0:
      urlist = eval(url)
      listed.append(urlist)
  return listed



#supporting functions: #printfancy
#print the contents of a list with sequence numbers and spaced out:
def printfancy(listo):
  seq = 1 
  if len(listo) > 0:
    for url in listo:
      print "(%3d)\t%s\n" % (seq, url)
      seq += 1
  else:
    print "empty list"


def printdict(dicto):
  ix = 1
  for clef in dicto.keys():
    print "[%d]\t%s=%s" % (ix, clef, dicto[clef])
    ix += 1



#supporting functions: #tabify
#put the contents of a list in a flat, tabbed string and return it:
def tabify(listo):
  strung =""
  if len(listo) == 0:
    return "empty list"
  for url in listo:
    if len(url) < 8:
      strung += "%s\t\t" % (url)
    else:
      strung += "%s\t" % (url)
  return strung


#supporting functions: #stringify
#put the contents of a list in a string and return it:
def stringify(listo):
  strung =""
  if len(listo) ==0:
    return "empty list"
  for url in listo:
    strung += "%s" % (url)
  return strung


#print the contents of a list if the line contains 'mode':
def cherrypick(listin, mode):
  modals = []
  if mode == "all":
    for el in listin:
      modals.append(el.strip())
  else:
    for el in listin:
      if el.find(mode) != -1:
        modals.append(el.strip())
  return modals


#add the elements of fromlist to tolist if not already in there.
def addunique(target, source):
  for item in source:
    if not inlist(target, item):
      target.append(item)
  return target


#return true if item is in the list else false
def inlist(tolist, elem):
  for one in tolist:
    if elem == one:
      return True
  return False


#if onestring contains any element from alist return true else false.
def findall(onestring, alist):
  for lis in alist:
    if onestring.find(lis) != -1:
      return True
  return False


#to flush the data: mode is "w" or "a":
def savefile(alist, filesaver, mode):
  try:
    ofp = open(filesaver, mode)
    for burl in alist:
      ofp.write("%s\n" % (burl))
    ofp.close()
  except IOError:
      #sys.stderr.write("savefile %s won't open\n" % (filesaver))
      raise

#to flush the data: mode is "w" or "a":
def savefilestring(astring, filesaver, mode):
  try:
    ofp = open(filesaver, mode)
    ofp.write("%s\n" % (astring))
    ofp.close()
  except IOError:
      #sys.stderr.write("savefile %s won't open\n" % (filesaver))
      raise

#Save dict as keys=values in the given file:
def savedict(adict, filesaver, mode):
  try:
    ofp = open(filesaver, mode)
    for burl in adict.keys():
      ofp.write("%s=%s\n" % (burl, adict[burl]))
  except IOError:
    raise


#open the given file and return its contents as a list.
def filegetter(filename):
  conts = []
  try:
    fp = open(filename)
    conts = fp.readlines()
    return conts
  except:
    return []


#strip newline from all lines in a list.
def stripper(inlist):
  unlist = []
  for el in inlist:
    unlist.append(el.strip())
  return unlist


#edit string takes a full url and appends or substitutes 
#the html filename supplied.
def editstring(atup, after):
  urlout = ""; threetup = []
  if atup.endswith("/"):
    urlout = atup + after
    return urlout

  threetup = atup.rpartition("/")
  urlout = atup.replace(threetup[2], after)
  return urlout


def visit(arg, dirname, names):
  allfiles = []
  for name in names:
    subname = os.path.join(dirname, name)
    if not os.path.isdir(subname):
      print '%s' % (subname)
    # Do not recurse into .svn directory
    if '.svn' in names:
      names.remove('.svn')



#get config file and import keyword arguments.
def getconfig(configfile):
  kv = {}
  arglines = filegetter(configfile)
  for argl in arglines:
    argyle = argl.strip()
    if argyle == "" or argyle.startswith("#"):
      continue
    kwargs = argyle.split("=")
    kv[kwargs[0]] = kwargs[1]
  return kv


#=========================================================================
if __name__ == "__main__":
  os.path.walk('d:\www.lan-opc.org.uk', visit, '')
  #fn = "parishesNov2010"
  #pars = filegetter(fn)
  #for par in pars: print par

'''
	print "helpers.py: Import functions as helpers:"
	print "\t functions include:\n\
\turlpeel:\tget an url and return contents as a list.\n\
\tsavefile:\tsave the list in the given file: mode is \"w\" or \"a\".\n\
\tfilegetter:\topen file and return contents as a list.\n\
\taddunique:\tcreate list as a set of unique elements.\n\
\tinlist: \tcheck if the given element is in the list.\n\
\tprintlist:\tprint out the given list.\n\
\tprintfancy:\tprint the given list with sequence numbers, spaced out.\n\
\tcherrypick:\tprint the list if the line contains 'mode'.\n\
\lstrip: \tstrip newline from all lines in a list."
'''
