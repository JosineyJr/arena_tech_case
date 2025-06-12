# Backend Developer Challenge – IP Location API

Your task is to implement a REST API using TypeScript that resolves a user's IP to a location by looking it up in a dataset.

---

## 🧠 Objective

Implement an endpoint `/ip/location` that accepts a query parameter called `ip` and returns the corresponding geographic location if found in the dataset.

---

## 📄 Dataset

You will receive a CSV file with **2,979,950 rows** ( `330Mb` ). Each row describes an IP range and its corresponding location.

- The file **does not** include column names.
- The IP ranges are given using **numeric representations of IPs**.
- Each row is structured as:

| Column | Description                                                       |
| ------ | ----------------------------------------------------------------- |
| 1      | Lower IP ID (inclusive)                                           |
| 2      | Upper IP ID (inclusive)                                           |
| 3      | Country Code                                                      |
| 4      | Country Name                                                      |
| 5      | State/Region                                                      |
| 6      | City                                                              |
| 7–10   | You may ignore these (latitude, longitude, postal code, timezone) |

---

## 🔢 IP to ID Conversion

> The provided CSV dataset is using an numeric ID to identify the IP

To convert a standard IPv4 address to an IP ID, use this formula:

```ts
function ipToId(ip: string): number {
    const ipParts = ip.split('.');
    let ipId = 0;
    for (let index = 0; index < ipParts.length; index++) {
      const element = parseInt(ipParts[index]);
      if (index === 0) {
        ipId += 16777216 * element;
      }
      if (index === 1) {
        ipId += 65536 * element;
      }
      if (index === 2) {
        ipId += 256 * element;
      }
      if (index === 3) {
        ipId += element;
      }
    }
    return ipId;
}
````

---

## 🌐 API Specification

### GET `/ip/location`

**Query Parameter**:

- `ip` (string): The IPv4 address to look up.
    

**Behavior**:

- Convert the IP to its corresponding IP ID using the `ipToId` function.
    
- Search for a row where:
    
    ```
    lower_ip_id <= ip_id <= upper_ip_id
    ```
    
- If no match is found, return:
    
```http
    Status: 404 Not Found
```
    
- If a match is found, return:
    
    ```json
    Status: 200 OK
    {
      "country": "Country Name",
      "countryCode": "Country Code",
      "city": "City"
    }
    ```
    

---

## 🧪 Load Testing with K6

Include a file `test.js` with the following content to validate the endpoint:

Run with:

```bash
yarn test:performance
```

> This test will validate a few valid and other invalid IP addresses and the minimal response time.
---

## 🧪 Starter API Example (Always returns 404)

Here’s a simple TypeScript/Restify starter that always returns 404, Feel free to use the framework you prefer:

### `index.ts`

```ts
import http, { IncomingMessage, ServerResponse } from 'http';
import { parse } from 'url';

const server = http.createServer((req: IncomingMessage, res: ServerResponse) => {
  const parsedUrl = parse(req.url || '', true);

  if (req.method === 'GET' && parsedUrl.pathname === '/ip/location') {
    res.statusCode = 404;
    res.setHeader('Content-Type', 'application/json');
    res.end(JSON.stringify({ error: 'Not Found' }));
    return;
  }

  // Default 404 response
  res.statusCode = 404;
  res.setHeader('Content-Type', 'application/json');
  res.end(JSON.stringify({ error: 'Not Found' }));
});

const PORT = 3000;
server.listen(PORT, () => {
  console.log(`Server listening at http://localhost:${PORT}`);
});

```

To run:

```bash
yarn add ts-node -D
yarn ts-node index.ts
```

---

## 🏗️ Evaluation Criteria

You will be evaluated on:

### Required

- **Correctness**: The API behaves according to the spec.
- **Architecture**: Clarity and organization of your code.
- **Code Quality**: Readability, modularity, maintainability.

### Desired

- **Performance**: Return the response in less than `100ms`
- **No third-party libraries**: If possible, don't use any third-party library except for the REST operations and/or automated tests.
- **Deal with concurrency**: It could be good if the API responds to more than `100` concurrent users.
- **Unit tests**: You could use [Jest](https://jestjs.io) or any other test library you prefer to implement unit tests.
- **Demonstrate your skills**: Put the most of your skills you can on this project to help us understand how you work and how you can help us.

---

## 📦 Submission

Please include:

- Your codebase in a zip file without the folder `node_modules`
- Clear instructions to install dependencies and run the server and tests.
- Your architectural decisions are in a `README.md` or similar file.
- Remove the file `IP2LOCATION-LITE-DB11.CSV` for the final zip.

