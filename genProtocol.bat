pushd protocol
..\tools\ProtocolGen.exe -p "TypeDefine.xml"			-b Buildin_Types.xml.cfg -s go -l utf-8 -m packet
..\tools\ProtocolGen.exe -p "TypeDefineCactus.xml"	    -b Buildin_Types.xml.cfg -s go -l utf-8 -m packet
popd
