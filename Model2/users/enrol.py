import sys, helpers
from Crypto.PublicKey import RSA
from Crypto.Signature import PKCS1_v1_5
from Crypto.Hash import SHA256
from Crypto import Random

#Methods for enrol.py:
#getKeys
#createKeys

#(args, priv, pub) = meth.createKeys(args)
def createKeys(args):
  rng = Random.new().read
  keypair = RSA.generate(1024, rng)
  privs = keypair.exportKey("PEM")
  priv = RSA.importKey(privs)
  helpers.savefilestring(privs, args['privkey'], 'w')
  pubs = keypair.publickey().exportKey("PEM")
  pub = RSA.importKey(pubs)
  helpers.savefilestring(pubs, args['pubkey'], 'w')
  return (privs, pubs)

#(priv, pub) = meth.getKeys(args)
def getKeys(args):
  privl = helpers.filegetter(args['privkey'])
  privs = helpers.stringify(privl)
  priv = RSA.importKey(privs)
  publ = helpers.filegetter(args['pubkey'])
  pubs = helpers.stringify(publ)
  pub = RSA.importKey(pubs)
  return (priv, pub)

#enrol.py <username>:
#Create a private and public key pair for the named user
#Save them in <user>.priv and <user>.pub files.

if __name__ == "__main__":
  args = {}

  try:
    user = sys.argv[1]
  except:
    sys.exit("Usage: python enrol.py <username>")

  priv = "%s.priv" % (user)
  pub = "%s.pub" % (user)
  args['privkey'] = priv
  args['pubkey'] = pub

  (privkey, pubkey) = createKeys(args)
  print privkey
  print pubkey

