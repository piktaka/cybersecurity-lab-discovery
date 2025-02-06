## install guacamol

#guacamole
cd /opt
apt install nginx -y
apt-get install tomcat9 tomcat9-admin tomcat9-common tomcat9-user -y
apt-get install -y libcairo2-dev libjpeg-turbo8-dev libpng-dev libtool-bin libossp-uuid-dev libavcodec-dev libavutil-dev libswscale-dev freerdp2-dev libpango1.0-dev libssh2-1-dev libvncserver-dev libtelnet-dev libssl-dev libvorbis-dev libwebp-dev
apt install gridsite-clients -y
wget https://archive.apache.org/dist/guacamole/1.5.0/source/guacamole-server-1.5.0.tar.gz
tar -xvf guacamole-server-1.5.0.tar.gz
cd guacamole-server-1.5.0
./configure --with-init-dir=/etc/init.d
make && make install
ldconfig
systemctl daemon-reload; systemctl enable guacd
mkdir /etc/guacamole
wget https://archive.apache.org/dist/guacamole/1.5.0/binary/guacamole-1.5.0.war -O /etc/guacamole/guacamole.war
ln -s /etc/guacamole/guacamole.war /var/lib/tomcat9/webapps/
#rdp
apt-get install xrdp -y
ufw allow 3389/tcp
#lightweight desktop env
DEBIAN_FRONTEND=noninteractive apt install xfce4 xfce4-goodies -y 
#to execute shotcuts
DEBIAN_FRONTEND=noninteractive apt-get -y install thunar 
sed -i "s/LogFile.*/LogFile=\/usr\/support\/xrdp-sesman.log/g" /etc/xrdp/sesman.ini
#disable rdp encryption
sed -i "s/crypt_level.*/crypt_level=none/g" /etc/xrdp/xrdp.ini
#for calculating timespent
DEBIAN_FRONTEND=noninteractive apt-get install -y xprintidle
sed -i "8s/3389/tcp:\/\/3389/g" /etc/xrdp/xrdp.ini
#change this otherwize xrdp will fail
sed -i 's/allowed_users=console/allowed_users=anybody/' /etc/X11/Xwrapper.config
#fix networking problem at startup
echo "auto eth0
iface eth0 inet dhcp" > /etc/network/interfaces
#xrdp-sesman failing to start
sed -i 's/forking/simple/g' /lib/systemd/system/xrdp-sesman.service
cp /lib/systemd/system/xrdp-sesman.service /etc/systemd/system/xrdp-sesman.service
sed -i '/SESMAN_OPTIONS/d' /etc/default/xrdp; sed -i '6i SESMAN_OPTIONS="-n"' /etc/default/xrdp
systemctl daemon-reload; systemctl restart xrdp-sesman xrdp



## install docker

apt update
apt install -y ca-certificates curl
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null
apt update
export DEBIAN_FRONTEND=noninteractive
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin




## installing burpsuite
INSTALL_DIR="/opt/BurpSuiteCommunity"
mkdir -p $INSTALL_DIR
cd $INSTALL_DIR
REFERER_URL="https://portswigger.net/burp/communitydownload"
OUTPUT_FILE="burpsuite_installer.sh"
# wget --referer="$REFERER_URL" https://portswigger-cdn.net/burp/releases/download?product=community&version=2024.11.2&type=Linux
wget --referer="$REFERER_URL" -O "$OUTPUT_FILE" "https://portswigger-cdn.net/burp/releases/download?product=community&version=2024.11.2&type=Linux"

chmod +x $OUTPUT_FILE
$INSTALL_DIR/burpsuite_installer.sh -q -dir "$INSTALL_DIR"

# cp $INSTALL_DIR/'Burp Suite Community Edition.desktop' /home/ubuntu/'Burp Suite'.desktop

# chmod x+ /home/ubuntu/'Burp Suite'.desktop


## installing firefox and hping3 and iptables

apt update && apt install -y firefox hping3 