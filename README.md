# genkit-go
Agentic patterns with genkit for go

Change the below parameters in the model for customized results from supported models.

Temperature

LLMs are fundamentally token-predicting machines. For a given sequence of tokens (such as the prompt) an LLM predicts, for each token in its vocabulary, the likelihood that the token comes next in the sequence. The temperature is a scaling factor by which these predictions are divided before being normalized to a probability between 0 and 1.

Low temperature values—between 0.0 and 1.0—amplify the difference in likelihoods between tokens, with the result that the model will be even less likely to produce a token it already evaluated to be unlikely. This is often perceived as output that is less creative. Although 0.0 is technically not a valid value, many models treat it as indicating that the model should behave deterministically, and to only consider the single most likely token.

High temperature values—those greater than 1.0—compress the differences in likelihoods between tokens, with the result that the model becomes more likely to produce tokens it had previously evaluated to be unlikely. This is often perceived as output that is more creative. Some model APIs impose a maximum temperature, often 2.0.

TopP

Top-p is a value between 0.0 and 1.0 that controls the number of possible tokens you want the model to consider, by specifying the cumulative probability of the tokens. For example, a value of 1.0 means to consider every possible token (but still take into account the probability of each token). A value of 0.4 means to only consider the most likely tokens, whose probabilities add up to 0.4, and to exclude the remaining tokens from consideration.

TopK

Top-k is an integer value that also controls the number of possible tokens you want the model to consider, but this time by explicitly specifying the maximum number of tokens. Specifying a value of 1 means that the model should behave deterministically.

