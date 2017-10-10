#/bin/sh




# build imagemagick
apt-get install -y gcc make git libwebp-dev libpng-dev libzopfli-dev libjpeg-dev curl wget bash tar xz-utils
IM_VER=ImageMagick-6.9.5-10
cd /tmp
wget http://www.imagemagick.org/download/releases/${IM_VER}.tar.xz
tar Jxvf ${IM_VER}.tar.xz
cd ${IM_VER}

./configure \
  --prefix=/usr \
  --enable-static \
  --without-lzma \
  --without-jbig \
  --disable-openmp

make && make install


# build pngquant
apt-get install -y gcc make git libpng-dev curl wget bash
PNGQUANT_VER=2.9.1
cd /tmp
git clone --recursive https://github.com/pornel/pngquant.git
cd pngquant
git checkout ${PNGQUANT_VER}
make && make install


# build optipng
apt-get install -y gcc make git libpng-dev curl wget bash
OPTIPNG_VER=0.7.6
cd /tmp
wget https://sourceforge.net/projects/optipng/files/OptiPNG/optipng-${OPTIPNG_VER}/optipng-${OPTIPNG_VER}.tar.gz
tar xvf optipng-${OPTIPNG_VER}.tar.gz
cd optipng-${OPTIPNG_VER}
./configure
make && make install


# build zopfli
apt-get install -y gcc make git libpng-dev libzopfli-dev curl wget bash
ZOPFLI_VER=zopfli-1.0.1
cd /tmp
git clone --recursive https://github.com/google/zopfli.git
cd zopfli
git checkout ${ZOPFLI_VER}
make && make zopflipng
mv zopfli /usr/local/bin
mv zopflipng /usr/local/bin

# build jpeg-archive
MOZJPEG_VER=3.2
JPEG_ARCHIVE_VER=2.1.1
# mozjpeg
apt-get install -y autoconf automake nasm libtool libjpeg-dev
cd /tmp
git clone https://github.com/mozilla/mozjpeg.git
cd mozjpeg
git checkout v${MOZJPEG_VER}
autoreconf -fiv
./configure --with-jpeg8
make
make install

# jpeg-archive
cd /tmp
git clone https://github.com/danielgtaylor/jpeg-archive.git
cd jpeg-archive
git checkout ${JPEG_ARCHIVE_VER}
make
make install


# go build
cd /go/src/github.com/kiwamunet/image-optim
go build .
mv image-optim /usr/local/bin
