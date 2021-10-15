import { newE2EPage } from '@stencil/core/testing';

describe('identifo-form', () => {
  it('renders', async () => {
    const page = await newE2EPage();
    await page.setContent('<identifo-form app-id="c3vqpr2se6nhtg9v1nu0" url="https://identifo.organaza.ru/" theme="dark"></identifo-form>');
    const element = await page.find('identifo-form');
    await page.compareScreenshot('My Component (...is beautiful. Look at it!)', { fullPage: false });
  });
  // it('renders changes to the name data', async () => {
  //   const page = await newE2EPage();
  //   await page.setContent('<identifo-form></identifo-form>');
  //   const component = await page.find('identifo-form');
  //   const element = await page.find('identifo-form >>> div');
  //   expect(element.textContent).toEqual(`Hello, World! I'm `);
  //   component.setProperty('first', 'James');
  //   await page.waitForChanges();
  //   expect(element.textContent).toEqual(`Hello, World! I'm James`);
  //   component.setProperty('last', 'Quincy');
  //   await page.waitForChanges();
  //   expect(element.textContent).toEqual(`Hello, World! I'm James Quincy`);
  //   component.setProperty('middle', 'Earl');
  //   await page.waitForChanges();
  //   expect(element.textContent).toEqual(`Hello, World! I'm James Earl Quincy`);
  // });
});
