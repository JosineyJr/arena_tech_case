import assert from 'assert';
import { beforeEach, afterEach, describe, it } from 'node:test';
import { start, stop } from '../server.ts';

describe('integration tests', () => {
  const PORT = 3999;
  const PATH = '/ip/location';
  beforeEach(() => {
    start(PORT);
  });

  afterEach(() => {
    stop();
  });

  it('should return 404 when the ip didn\'t exists', async () => {
    const result = await fetch(`http://localhost:${PORT}/${PATH}?ip=255.255.255.255`);
    assert.strictEqual(result.status, 404, 'should return the status code 404');
  })
})

