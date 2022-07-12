import {
  ApolloClient,
  createHttpLink,
  fromPromise,
  InMemoryCache,
} from "@apollo/client/core";
import { onError } from "@apollo/client/link/error";
import { useErrorsStore } from "../stores/errors";
import { setContext } from "@apollo/client/link/context";
import { useAuthStore } from "../stores/auth";

const httpLink = createHttpLink({
  uri: "http://localhost:8080/query",
});

const errorHandler = onError(({ graphQLErrors, operation, forward }) => {
  if (graphQLErrors) {
    for (let err of graphQLErrors) {
      switch (err.message) {
        case "Token is expired":
          const auth = useAuthStore();
          // regenerate token
          console.log("Token is expired, regenerating from refresh token");
          return fromPromise(
            auth.refresh().catch((error) => {
              return false;
            })
          )
            .filter((value) => {
              console.log("Refreshed, value=" + value);
              return false;
            })
            .flatMap((accessToken) => {
              const oldHeaders = operation.getContext().headers;
              // modify the operation context with a new token
              operation.setContext({
                headers: {
                  ...oldHeaders,
                  authorization: `Bearer ${accessToken}`,
                },
              });

              // retry the request, returning the new observable
              return forward(operation);
            });
      }
    }
    useErrorsStore().$state = {
      message: graphQLErrors[0].message,
      category: graphQLErrors[0].extensions?.category as any,
      fields: graphQLErrors[0].extensions?.validation ?? ({ input: {} } as any),
    };
  }
});

const authLink = setContext((operation, prevContext) => {
  const authStore = useAuthStore();

  const res = {
    headers: {
      ...prevContext.headers,
    },
  } as any;
  console.log(`Operation name: ${operation.variables} : ${prevContext}`);
  console.log(operation.extensions);
  if (authStore.accessToken && operation.operationName !== "renewToken") {
    res.headers.authorization = "Bearer " + authStore.accessToken;
  }
  return res;
});

// Cache implementation
const cache = new InMemoryCache();

// Create the apollo client
const apolloClient = new ApolloClient({
  link: authLink.concat(errorHandler.concat(httpLink)),
  cache,
});

export default apolloClient;
