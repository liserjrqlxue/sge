# sge
qsub -cwd -l p=n,vf=mG -binding linear:n -P project -q bc.q job.sh
