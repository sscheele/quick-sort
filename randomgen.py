from random import choice
import string

with open("test.dat", 'w') as f:
	for _ in range(1000):
		f.write(''.join(choice(string.ascii_uppercase) for _ in range(30)) + "\n")
