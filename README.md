
# Recovery tool

A grep-like tool that is effective on large binary files, written in Go.

It is useful for data recovery on defect disk images, recovery of accidentally deleted files, etc.
This tool can search for some known text passages (or byte arrays) on a whole disk image,
so it might be possible to recover some lost data, e.g. ASCII text files.

## Example uses

Create a disk image (same as `dd if=/dev/sda1 of=disk.img -oflag=sync -bs=100M`).
It is a little bit faster than `dd` as `recoverytool` uses separate threads for reading and writing.
It also prints a status report and the MD5 checksum.

	sudo ./recoverytool -cp img.iso /dev/sda1

Search for byte patterns in a file image.
Also works on a whole disk.

	./recoverytool -m "Master's Thesis" -m "Tax report" -m "Dear Mr. Santa" disk.img > report.txt

It takes some time to run on large file images.
If a pattern was found, it prints out the address in hex and some lines of context.
If you stream the output to a file, like in the example above, 
you can let it run and analyze the result later.

Here is an example output of a simple search:

```
./recoverytool -m "debug" ../comcat/comcat 
> Searching ../comcat/comcat from 0 to 175e8: 95720 Bytes (95.7kB)
> Match for debug at adress 16c0d
> 16a80 .restore_stdin._ITM_registerTMCloneTable.__cxa_finalize@@GLIBC_2.2.5.__ctype_b_loc@@GLIBC_2.3.stderr@@GLIBC_2.2.5...symtab..strt
> 16b00 ab..shstrtab..interp..note.gnu.property..note.gnu.build-id..note.ABI-tag..gnu.hash..dynsym..dynstr..gnu.version..gnu.version_r..
> 16b80 rela.dyn..rela.plt..init..plt.got..plt.sec..text..fini..rodata..eh_frame_hdr..eh_frame..init_array..fini_array..dynamic..data..b
> 16c00 ss..comment..debug_aranges..debug_info..debug_abbrev..debug_line..debug_str..debug_ranges..debug_macro..........................
> 16c80 ........................................................................................................#...............8.......
> 16d00 8....... ...............................6...............X.......X.......$...............................I...............|.......
> 16d80 |....... ...............................W......o........................0...............................a.......................
> Match for debug at adress 16c1c
> 16a80 .restore_stdin._ITM_registerTMCloneTable.__cxa_finalize@@GLIBC_2.2.5.__ctype_b_loc@@GLIBC_2.3.stderr@@GLIBC_2.2.5...symtab..strt
> 16b00 ab..shstrtab..interp..note.gnu.property..note.gnu.build-id..note.ABI-tag..gnu.hash..dynsym..dynstr..gnu.version..gnu.version_r..
> 16b80 rela.dyn..rela.plt..init..plt.got..plt.sec..text..fini..rodata..eh_frame_hdr..eh_frame..init_array..fini_array..dynamic..data..b
> 16c00 ss..comment..debug_aranges..debug_info..debug_abbrev..debug_line..debug_str..debug_ranges..debug_macro..........................
> 16c80 ........................................................................................................#...............8.......
> 16d00 8....... ...............................6...............X.......X.......$...............................I...............|.......
> 16d80 |....... ...............................W......o........................0...............................a.......................

...
```

With a report, you can analyze what is going on around address 16c0d in the disk image:

	./recoverytool -s 16c0d -n 100 -p disk.img

The length of a printed line is 128, by default.
In the hex dump produced by `-p`, ASCII characters between #32 (`' '`, space) and #126 (`~`) are printed as they are.
The characters `\t`, `\r`, `\n`, `\f` are printed as spaces, and all other characters are printed as dots (`.`).
There is no unicode support.

You can also calculate the MD5 checksum of a file:

	./recoverytool -md5 disk.img

# Complete Usage and Command Line Arguments

```
Usage: recoverytool <options> <action> <file name>


Options:

  -b <no of bytes>
        Buffer size used. Should be plenty to avoid kernel overhead. (default 134217728)
  -s <offset>
        Begin at this byte offset (inclusively) in the source file.
  -e <offset>
        Read up to this byte offset (exclusively) in the file.
        A negative value allows to specify an ofset from the end of the file.
        A value of 0 indicates the end of the file.
  -l <number of bytes>
        Length of a "line" for report (default 128).
  -n <number of lines>
        Number of "lines" to process, as an alternative to -e flag.
  -la <number of lines>
        "Lines" of context to report after a match (default 3).
  -lb <number of lines>
        "Lines" of context to report before a match (default 3).
  -q    
        Do not print status updates.
  -v    
        Be more verbose.
  -d    
        Print debug output.
  -h
        Print a help message and exit.

        
The following actions can be chosen:

  -cp <target filename>
        Similar to dd, copy file contents to another file.
  -m <string>
        Search for a pattern. 
        Can be specified multiple times to search for several patterns in one go.
  -md5
        Calculate the md5 checksum.
  -p    
        Print out a hex dump.
```
