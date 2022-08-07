TODO
## Background
This project simulates a virtual world with countries, resources, and actions that these countries can perform with these resources. It is a demonstration of artificial intelligence through the use of a forward searching, utility driven, depth-bound, anytime search.

* Forward searching: The search will move "forward" only, and will not search through states behind it.

* Utility driven: This means that it will search for a best path of actions to take for the "self" country, using an expected utility score.

* Depth-bound: The depth-first nature of the search will be stopped at a certain depth in order to evaluate the best state based on the expected utilities.

* Anytime: This means the search won't stop once it hits the depth bound for the first time. It will continue search other paths until its priority queue has filled up.

This search can also be categorized as the following types of search:
* Best first search, since it chooses the state of the frontier to move based on the best nodes first
* Beam search, since it will include the top X amount of "best" actions based on the parameters passed in
* Heuristic search, since it will store the actions/states based on their expected utility in a way that allows it to pop the best node off of the priority queue first.

## Additional Details
The project will construct the "real" world based on the inputs passed in, and then it will begin to simulate all of the possible actions that could be taken by or with the "self" country that is passed in. After the frontier (i.e. the priority queue) has filled up, then it will backup the best states from the frontier, advance the "real" world with the best state, and then flush the priority queue. The queue will then be re-initialized with the previous best states as a way to continue the search with these values, only now one more state ahead.

Since we are using a depth-bounded search, this means that each round, the depth we search is also increased by one so that our depth bound maintains a consistent distance from our "real" state.

The sequences of actions that we simulate are called the "schedule", and they are stored in a priority queue using an Expected Utility evaluation. For further details on the calculations, there is a Calculations section further down this README.

## Running
With Golang 1.18 installed and your GOROOT and GOPATH set correctly, you should be able to run the following commands:

To build and run the code:
` $ go build main.go `
` $ <executable> <self country name> <resources csv file> <initial state csv file> <schedule output filename> <max top states> <max depth each round> <max frontier size> <number of rounds>`

For instance, an example run on my windows computer is:
` $ main.exe "Gondor", "..\inputs\resources.csv", "..\inputs\countries_balanced.csv", "..\test_outputs\output.txt", 4, 5, 100, 10`

Here is a break down of the input parameters:
* \<executable>: This is the executable created by the `$ go build main.go` script
* \<self country name>: This is the name of the "self" country for this run
* \<resources csv file>: This is the file containing the resource names and weights
* \<initial state csv file>: This is the file containing the countries alongside their initial state (amount of each resource)
* \<schedule output filename>: This is the name of the output file for the resultant schedules
* \<max top states>: This is the max number of "best" states that are used after each round as a part of the beam search
* \<max depth each round>: This is the maximum depth to search each round before spreading out the search
* \<max frontier size>: This is the maximum size of the frontier (in this scenario, that means the priority queue)
* \<number of rounds>: This is the number of rounds to evaluate

To run the test code:
`$ go test`


## Calculations
There are several calculations that are worth covering here: State Quality, Undiscounted Reward, Discounted Reward, Success Probability, and Expected Utility

### State Quality
A country's State Quality is calculated from a sum-of-parts state quality for the different aspects of the resources I am using.
* Food Quality = (WeightedFoodValue + WeightedFoodWasteValue) / Population
  * Justification: 
* Housing Quality = (WeightedHousingValue + WeightedHousingWasteValue) / (Population * 0.4)
  * Justification: I am only using 40% of the population value since multiple people will live in one house. According to the U.S. Census Bureau, that number was 2.60 people per household from 2016-2020 https://www.census.gov/quickfacts/fact/table/US/HCN010217 (use the dropdown to select Families & Living Arrangements) 
* Electronics Quality = WeightedElectronicsValue + WeightedElectronicsWasteValue
  * Justification: I am not dividing this by any value since I'm considering it more of a generic "technology" score of a country.
* Metals Quality = WeightedMetallicAlloyValue + WeightedMetallicAlloyWasteValue
  * Justification: I am not dividing this by any value because they will be used for a variety of other calculations
* Land Quality = (WeightFarmValue^2 + WeightFarmWasteValue) / AvailableLand
  * Justification: Since Farmland takes up a lot of space, I am including this in the Land Quality by squaring the weighted value for farmland. 
* Military Quality = (WeightedMilitaryWeight + WeightedMilitaryWasteWeight) / (Population * 0.2)
  * Justification: For the countries I have passed in based on the world Lord of the Rings by JRR Tolkien, the military would be important for defense. However it is a costly resource to create, so I am only dividing by 20% of the population.

I really wanted to generate the partial qualities in isolation based on the most meaningful measure for that quality. However, the side effect of that is that these values need to be rebalanced so that each produces a "fair" part of the overall state quality. 
My fictional countries and the world I was trying to simulate is based on the world of the Lord of the Rings by JRR Tolkien, this means that the different partial quality weights would need once more rebalancing according to this test data. In essence, this means that there are 3 layers of rebalancing going on:
1. Each resource has a weight for how it gets used
2. Each Partial Quality has a formula to be meaningfully calculated based on the weighted components
3. The State Quality then uses weights on each of these Partial Qualities to rebalance them against one another.
I have added additional weights for each partial quality in order to weigh them against each other. \
Therefore, the State Quality formla is:
> State Quality = (Food Quality * 1.00) + (Housing Quality * 0.3) + (Electronics Quality * 0.15) + (Metals Quality * 0.15) + (Land Quality * 0.2) + (Military Quality * 0.5)

Under these terms that I have defined, I am considering Food to be a critical resource, followed by military. This felt intuitive for a medieval society. Next I have housing followed by land use as the next most valuable resources. The electronics (which I consider to be a more generalized "technology" for this world) and metals are weighted equally at the bottom, because while both have value, they aren't as valuable as the other resources.

### Undiscounted Reward
A country's Undiscounted Reward is the difference between that country's current State Quality, and the resultant State Quality after a schedule is executed

### Discounted Reward
A country's Discounted Reward is the Undiscounted Reward multiplied by a discount value. The Discount Value is _gamma^n_, where _gamma_ is a configurable constant, and _n_ is the distance from the current state of the world. This calculation means that an updated state will be more valuable the fewer steps you have to take to get there. 

### Success Probability
A schedule's Success Probability is calculated by taking every country that's a part of the schedule, calculating the discounted rewards, passing them through a Logistic Regression function, and lastly multiplying them together. The product is the probability that the schedule will be accepted

### Expected Utility
The Expected Utility is calculated for the "self" country by taking its Discounted Reward, multiplying it by the Success Probability of the schedule, and then adding this to the Failure Probability (1 - Success Probability) multiplied by the Failure Cost (in this case a configurable constant).

### Additional Resources and Justification:
I am converting available land into farmland in a 5-to-1 ratio, and then converting 1 farm to 1 food. This is justified by the following source I found claiming that 5 or 6 acres of land is needed for someone to comfortable produce their own food. https://permaculturism.com/how-much-land-does-it-take-to-feed-one-person/#:~:text=A%20person%20needs%20about%205,comfortably%2C%20producing%20their%20own%20food.

In order to produce the Military Resource, I am converting food, population, housing, and metallic alloys. Food is in a  3-to-1 ratio according to this stackexchange post https://worldbuilding.stackexchange.com/questions/13081/how-many-resources-are-required-to-feed-1-million-soldiers, plus military will need places to sleep as well as metal for weapons, so these were factored in using 1-to-1 and 2-to-1 ratios respectively.


## Repo Structure
* _cmd_: This folder has both the main.go and main_test.go files. The executable will be generated here from the commands above.
* _Manager_: This contains the game manager that monitors each country for proposed schedules and determines the best overall schedule.
* _Scheduler_: Holds the actual scheduler functions and maintains all countries, resources, and actions that can be used in a schedule.
* _Simulator_: Holds the simulator that will duplicate the world state, simulate actions that can be taken, and then search for the best action. It will then pass this back to the main process that will tell the scheduler to enact the first action in this schedule.
* _Resources_: Holds the resource struct and functions.
* _Countries_: Holds the country struct and functions.
* _Quality_: Holds the state quality function along with the other reward calculations.
* _Transfers_: Holds the transfer action struct and its functions.
* _Transformations_: Holds the transformation action struct and its functions.
* _Util_: Holds utility functions and constants.
* _input_: Holds the input files that can be used (3 country files and 1 resources file).
* _test_outputs_: Holds the outputs generated by the program run and the test case run.

