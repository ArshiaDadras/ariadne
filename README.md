# Ariadne - A Map Matching Engine in Golang

## About Ariadne

In Greek mythology, Ariadne was a princess of Crete who played a crucial role in the myth of the Labyrinth. She is most famously known for aiding Theseus in navigating the Labyrinth by providing him with a thread to trace his path back out after slaying the Minotaur. This thread became a symbol of guidance and a means of finding one's way through complex situations.

To learn more about Ariadne's mythological significance, visit her [Wikipedia page](https://en.wikipedia.org/wiki/Ariadne).

## What is a Map Matching Engine?

A map matching engine is a computational tool used to align or "match" a sequence of geographic coordinates (typically obtained from GPS devices) to the most likely path on a map. This process is essential for applications such as navigation systems, route planning, and geographic data analysis, where accurate positioning relative to a road network or predefined paths is crucial.

## The Hidden Markov Model (HMM) Technique

The core of the Ariadne map matching engine leverages the Hidden Markov Model (HMM) technique, a probabilistic approach used to infer the most probable path a vehicle or individual has taken, given noisy or sparse data. The HMM technique involves:

1. **State Representation**: In the context of map matching, states represent the possible positions or paths on the map.
2. **Observation Model**: This model accounts for the GPS data (observations) and how they relate to the actual states on the map.
3. **Transition Model**: This defines the probabilities of moving from one state (location on the map) to another, considering factors like road networks and travel dynamics.

By applying the HMM approach, Ariadne effectively filters out noise and infers the most likely route that matches the observed GPS coordinates with the underlying map.

For a more detailed understanding of the HMM map matching technique, you can refer to the [Microsoft Research publication](https://www.microsoft.com/en-us/research/publication/hidden-markov-map-matching-noise-sparseness/).

## License

Ariadne is licensed under the [MIT License](LICENSE). See the LICENSE file for more details.
