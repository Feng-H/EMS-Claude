const puppeteer = require('/private/tmp/puppeteer-test/node_modules/puppeteer');

(async () => {
  const browser = await puppeteer.launch({ 
    headless: 'new', 
    args: ['--no-sandbox', '--disable-setuid-sandbox'] 
  });
  const page = await browser.newPage();
  await page.setUserAgent('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36');
  await page.setDefaultNavigationTimeout(60000); 
  
  console.log('开始全面体检：https://ems.317316.xyz');

  try {
    // 0. 诊断：页面源码检查
    console.log('正在探测页面结构...');
    await page.goto('https://ems.317316.xyz', { waitUntil: 'networkidle2' });
    const content = await page.content();
    console.log('页面标题:', await page.title());
    console.log('当前 URL:', page.url());
    
    if (content.includes('404') || content.includes('Not Found')) {
        console.log('错误：页面返回了 404');
    }

    // 1. 登录测试
    console.log('[1/4] 正在尝试登录流程...');
    if (!page.url().includes('login')) {
        await page.goto('https://ems.317316.xyz/login', { waitUntil: 'networkidle2' });
    }
    
    // 等待任何输入框出现
    await page.waitForSelector('input');
    const inputs = await page.$$('input');
    if (inputs.length >= 2) {
      await inputs[0].type('admin');
      await inputs[1].type('admin123');
    }
    
    const loginBtn = await page.evaluateHandle(() => {
      const btns = Array.from(document.querySelectorAll('button'));
      return btns.find(b => b.textContent.includes('登录') || b.className.includes('primary'));
    });
    if (loginBtn) await loginBtn.click();
    
    await page.waitForNavigation({ waitUntil: 'networkidle2' });
    console.log('登录成功！');

    // 2. 检查控制台
    console.log('[2/4] 检查工作台数据...');
    await page.goto('https://ems.317316.xyz/dashboard', { waitUntil: 'networkidle2' });
    const hasCards = await page.evaluate(() => document.querySelectorAll('.el-card').length > 0);
    console.log(hasCards ? '工作台数据加载正常' : '警告：工作台似乎没有数据');

    // 3. 检查飞书配置页面
    console.log('[3/4] 检查飞书配置页...');
    await page.goto('https://ems.317316.xyz/agent/assistant', { waitUntil: 'networkidle2' });
    
    // 点击绑定飞书
    await page.evaluate(() => {
      const items = Array.from(document.querySelectorAll('.nav-item'));
      const bindItem = items.find(i => i.textContent.includes('绑定飞书账号'));
      if (bindItem) bindItem.click();
    });
    await new Promise(r => setTimeout(r, 2000));

    // 检查表单是否出现
    const isFormVisible = await page.evaluate(() => !!document.querySelector('input[placeholder="cli_xxxxxxxx"]'));
    console.log(isFormVisible ? '飞书配置弹窗正常弹出' : '错误：飞书配置弹窗未显示');

    // 4. H5 模拟绑定流程测试
    console.log('[4/4] 模拟 H5 绑定链路...');
    const h5Btn = await page.evaluateHandle(() => {
      const btns = Array.from(document.querySelectorAll('button'));
      return btns.find(b => b.textContent.includes('直接前往 H5 模拟绑定'));
    });

    if (h5Btn) {
      const newPagePromise = new Promise(x => browser.once('targetcreated', target => x(target.page())));
      await h5Btn.click();
      const h5Page = await newPagePromise;
      console.log('H5 页面已打开:', h5Page.url());
      
      await h5Page.waitForNavigation({ waitUntil: 'networkidle2' }).catch(() => {});
      
      // 捕获 API 报错
      h5Page.on('response', async response => {
        if (response.url().includes('bind-lark') && !response.ok()) {
          const body = await response.json();
          console.log('!!! 核心错误捕获 !!!');
          console.log('URL:', response.url());
          console.log('Status:', response.status());
          console.log('Error:', JSON.stringify(body));
        }
      });

      console.log('尝试点击立即绑定...');
      await h5Page.evaluate(() => {
        const btns = Array.from(document.querySelectorAll('button'));
        const bindBtn = btns.find(b => b.textContent.includes('立即绑定'));
        if (bindBtn) bindBtn.click();
      });
      await new Promise(r => setTimeout(r, 3000));
    }

  } catch (err) {
    console.error('测试中断:', err.message);
  } finally {
    console.log('测试结束，正在关闭。');
    await browser.close();
  }
})();
