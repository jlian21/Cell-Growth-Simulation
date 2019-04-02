# Cell-Growth-Simulation

1. Download the file and unzip it to the go/src directory.
2. Change the directory of the simulator
3. Compile all the packages (go build) in the directory. 

On Mac, the command for the simulator follows the format:

./cgsimu	COMMAND		OPTION

There are three options for the COMMAND. 

#######
The FIRST option is OneCluster. 

Taking the command, the simulator reads the stored inputs in the OneClusterOutputs.txt, and generates a gif, named OneCluster.gif, that records the simulation of one cluster of cells according to the inputs. 

OPTION available for OneCluster are Voronoi and CountDensity, which are two ways of counting the surrounding density of a cell.
#######

#######
The SECOND option is AutoGenerate. 

Under this command, the simulator randomly generates inputs, stores as input.txt, and simulates the growth pattern of one cluster of cells according to the generated inputs, and the output gif is named AutoGenerate.gif. 

OPTION available for AutoGenerate are Voronoi and CountDensity
#######

#######
The THIRD option is TwoCluster. 

The simulator reads the stored inputs in the TwoClusterOutputs.txt, and generates a gif, named TwoClusterSS.gif, recording the simulation of the source and sink model according to the inputs.

NO OPTION available for TwoCluster, and the default counting method is CountDensity.
#######

+++++++
Tips:

For ./cgsimu OneCluster CountDensity, it usually takes 30-50 seconds for 300 generations.

If using Voronoi, the recommended number of generation is < 100 (for now).

Also, if (unfortunately) encountering index out of range, please type the command again since randomness can be nasty at times :)
+++++++
