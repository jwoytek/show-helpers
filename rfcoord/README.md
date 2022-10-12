# RF Coordination Stuff
## scanfixer2.py
This is a simple python script that will convert a scan output from a
Siglent spectrum analyzer into a format that can be imported into 
Shure WWB and used for RF coordination. 

To use: 
  1. Save one or more scans for the ranges of interest to the USB
stick on the Siglent. 
  2. Plug the stick into a computer with Python 3 installed.
  3. Run the command for each scan file, replacing the two filenames as needed: `./scanfixer2.py siglent_scan_file.csv file_for_wwb.csv`
    * NOTE: This will happily overwrite whatever you specify as the output file, so use caution when specifying the filenames!
  4. Import converted scans into WWB.

