## Configure file

#### TuXiv Directory tree

<img src="/Users/zhangcengguang/Library/Application Support/typora-user-images/image-20200905105043809.png" alt="image-20200905105043809" style="zoom:40%;" />

#### Model configuration

+ Start Scripts

  ~~~yaml
  start_cmd: python train.py 
  ~~~

#### Conda configuration

+ Export(https://web.archive.org/web/20170704223211/https://www.continuum.io/content/conda-data-science)

  ~~~shell
  conda env export > conf.yaml
  ~~~

+ Conf.yaml example

  ~~~yaml
  name: sensor_gateway
  channels:
    - defaults
  dependencies:
    - aiohttp=2.3.9=py36_0
    - async-timeout=2.0.0=py36hc3e01a3_0
    - certifi=2018.4.16=py36_0
    - chardet=3.0.4=py36h420ce6e_1
    - icc_rt=2017.0.4=h97af966_0
    - intel-openmp=2018.0.0=8
    - jinja2=2.10=py36h292fed1_0
    - markupsafe=1.0=py36h0e26971_1
    - mkl=2018.0.2=1
    - mkl_fft=1.0.1=py36h452e1ab_0
    - mkl_random=1.0.1=py36h9258bd6_0
    - multidict=3.3.2=py36h72bac45_0
    - numpy=1.14.2=py36h5c71026_1
    - pip=9.0.1=py36h226ae91_4
    - pymysql=0.7.11=py36hf59f3ba_0
    - python=3.6.4=h6538335_1
    - pytz=2018.3=py36_0
    - setuptools=38.4.0=py36_0
    - simplejson=3.14.0=py36hfa6e2cd_0
    - sqlalchemy=1.2.1=py36hfa6e2cd_0
    - vc=14=h0510ff6_3
    - vs2015_runtime=14.0.25123=3
    - wheel=0.30.0=py36h6c3ec14_1
    - wincertstore=0.2=py36h7fe50ca_0
    - yarl=0.14.2=py36h27d1bf2_0
  prefix: C:\ProgramData\Anaconda3\envs\sensor_gateway
  ~~~

+ Import

  ~~~shell
  conda env create -f conf.yaml
  ~~~

#### Slurm configuration

+ Export

  + copy from /etc/slurm.conf

+ Conf example

  ~~~shell
  #
  # Sample /etc/slurm.conf for dev[0-25].llnl.gov
  # Author: John Doe
  # Date: 11/06/2001
  #
  SlurmctldHost=dev0(12.34.56.78) # Primary server
  SlurmctldHost=dev1(12.34.56.79) # Backup server
  #
  AuthType=auth/munge
  Epilog=/usr/local/slurm/epilog
  Prolog=/usr/local/slurm/prolog
  FirstJobId=65536
  InactiveLimit=120
  JobCompType=jobcomp/filetxt
  JobCompLoc=/var/log/slurm/jobcomp
  KillWait=30
  MaxJobCount=10000
  MinJobAge=3600
  PluginDir=/usr/local/lib:/usr/local/slurm/lib
  ReturnToService=0
  SchedulerType=sched/backfill
  SlurmctldLogFile=/var/log/slurm/slurmctld.log
  SlurmdLogFile=/var/log/slurm/slurmd.log
  SlurmctldPort=7002
  SlurmdPort=7003
  SlurmdSpoolDir=/var/spool/slurmd.spool
  StateSaveLocation=/var/spool/slurm.state
  SwitchType=switch/none
  TmpFS=/tmp
  WaitTime=30
  JobCredentialPrivateKey=/usr/local/slurm/private.key
  JobCredentialPublicCertificate=/usr/local/slurm/public.cert
  #
  # Node Configurations
  #
  NodeName=DEFAULT CPUs=2 RealMemory=2000 TmpDisk=64000
  NodeName=DEFAULT State=UNKNOWN
  NodeName=dev[0-25] NodeAddr=edev[0-25] Weight=16
  # Update records for specific DOWN nodes
  DownNodes=dev20 State=DOWN Reason="power,ETA=Dec25"
  #
  # Partition Configurations
  #
  PartitionName=DEFAULT MaxTime=30 MaxNodes=10 State=UP
  PartitionName=debug Nodes=dev[0-8,18-25] Default=YES
  PartitionName=batch Nodes=dev[9-17] MinNodes=4
  PartitionName=long Nodes=dev[9-17] MaxTime=120 AllowGroups=admin
  ~~~

+ Import

  + Stop the Slurm daemons
  2. Modify the slurm.conf file appropriately
  3. Distribute the updated slurm.conf file to all nodes
  4. Restart the Slurm daemons

  

