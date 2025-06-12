import assert from 'assert';
import { describe, it } from 'node:test';
import { ipToID } from '../ip-calculator.ts';

describe('ipToID', () => {
  it('should return a number', () => {
    const result = ipToID('0.0.0.0');
    assert(typeof result === 'number', 'it should be a number');
  });

  it('when receive 8.8.8.8 should return 0', () => {
    const result = ipToID('8.8.8.8');
    assert.strictEqual(result, 134744072, 'it should be a number');
  });
});

