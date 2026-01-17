const bcrypt = require('bcryptjs');

// 默认密码
const password = 'password123';

// 生成hash
const hash = bcrypt.hashSync(password, 10);

console.log('========================================');
console.log('Password Hash Generator');
console.log('========================================');
console.log(`\nPassword: ${password}`);
console.log(`\nBcrypt Hash (cost=10):`);
console.log(hash);
console.log('\n========================================');
console.log('\nCopy this hash and use it in seed.sql');
console.log('Replace all instances of: $2a$10$YourHashedPasswordHere');
console.log('========================================\n');

// 生成多个hash供验证（应该都相同）
console.log('Verification hashes (should all match):');
for (let i = 0; i < 3; i++) {
  console.log(`  ${i + 1}. ${bcrypt.hashSync(password, 10)}`);
}

// 验证hash
console.log('\nVerification:');
console.log(`Verify password against hash: ${bcrypt.compareSync(password, hash)}`);
