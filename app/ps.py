import scipy.stats as st
import numpy as np
import cv2

def brighten(img, degree):
	hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
	h,s,v = cv2.split(hsv)
	lim = 255 - degree
	v = np.clip(v,0,lim)
	v += degree
	hsv = cv2.merge((h,s,v))
	res = cv2.cvtColor(hsv, cv2.COLOR_HSV2BGR)
	return res

def saturation(img, degree):
	hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV).astype("float64")
	(h, s, v) = cv2.split(hsv)
	s += degree
	s = np.clip(s,0,255)
	hsv = cv2.merge((h,s,v))
	res = cv2.cvtColor(hsv.astype("uint8"), cv2.COLOR_HSV2BGR)
	return res

def contrast(img, degree):
	lab = cv2.cvtColor(img,cv2.COLOR_BGR2LAB)
	l,a,b = cv2.split(lab)
	clahe = cv2.createCLAHE(clipLimit=3.0, tileGridSize=(8,8))
	cl = clahe.apply(l)
	limg = cv2.merge((cl,a,b))
	res = cv2.cvtColor(limg, cv2.COLOR_LAB2BGR)
	return res

def crop(img, lt, rb):
	return img[lt[1]:rb[1], lt[0]:rb[0]]

def rotateImg(img, degree):
	rows,cols = img.shape[0], img.shape[1]
	rot_mat = cv2.getRotationMatrix2D((cols/2,rows/2),degree,1.0)
	res = cv2.warpAffine(img, rot_mat, (cols,rows),flags=cv2.INTER_LINEAR)
	return res

def gaussianBlur(img, klen):
	sig = klen/7.
	interval = (2*sig+1.)/(klen)
	x = np.linspace(-sig-interval/2., sig+interval/2., klen+1)
	k1d = np.diff(st.norm.cdf(x))
	kraw = np.sqrt(np.outer(k1d,k1d))
	G = kraw/kraw.sum()
	return cv2.filter2D(img,-1,G)

def bilareralBlur(img, Blen, Glen):
	# HighPass = img.copy()
	HighPass = cv2.bilateralFilter(img,Blen/7,Blen,Blen)
	# HighPass = HighPass - img + 128
	# layer1 = gaussianBlur(HighPass,Glen)
	# layer2 = HighPass + img - 128
	# return layer1/2 + layer2/2
	return HighPass

def deNoise(img, klen):
	return cv2.fastNlMeansDenoisingColored(img,None,klen,klen,7,15)
