# Readme

Problem: https://pastebin.com/gevUN7Be  


**Dwarves Foundation**  
Email: team@dwarvesv.com   
Website: https://dwarves.foudation/  

## Specifications
* Language: **Golang**  
* Database: **LevelDB**   
* Algorithm: **Depth First Search (DFS)**  


## Setup and Run


## Features
We have came up with these techniques and optimization to boost up the solution:

**1.  Data Generation**  
Because the problem statement does not specify how the generated data should be, we are leveraging the waterfall model to avoid being cyclic while building the Graph. The waterfall model means water only flows from the high level to the lower level and not the other way around.

*For example* .   
The value of a vertex is an integer. We set a rule that there is only one-way edge connecting 2 edges with the direction from the high value to the lower value.

If the vertices (12) and (4) are connected, there is only 1 edge connecting (12) --> (4) and no the other ways around in any cases (say 4 --> 12 for example)

With the given vertices and edges generated leveraging this model, we don't need to worry that the DAG is being cyclic while constructing the DAG from the datastream.

This waterfall model is also applied for our topological order.

**2.  DAG Construction**  
We are using DFS to
* implement these functions
    1) Reach(vertex ID)
    2) ConditionalReach(vertex ID, flagCondition)
    3) List(vertex ID) + ConditionalList(vertex ID, flagCondition)
* ensure that the DAG is not being cyclic in the case the data does not follow the waterfall model.

We are also connecting disconnected components in the DAG in the waterfall topologial order to make sure it is well-connected.


**3.  Concurrency & Optimization**
* We do not load all vertices into memory, try to load one-by-one vertex to process.
* **Go Routines** with wait group to handle when process done. Use channel to receive data.
* When inserting a vertex to the map, we create a virtual map and update the virtual map into the database.
