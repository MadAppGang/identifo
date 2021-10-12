import { newSpecPage } from '@stencil/core/testing';
import { IdentifoForm } from './identifo-form';

describe('identifo-form', () => {
  it('renders', async () => {
    const { root } = await newSpecPage({
      components: [IdentifoForm],
      html: '<identifo-form app-id="c3vqpr2se6nhtg9v1nu0" url="https://identifo.organaza.ru/" theme="dark"></identifo-form>',
    });
    expect(root).toEqualHtml(`
        <identifo-form app-id="c3vqpr2se6nhtg9v1nu0" route="error" url="https://identifo.organaza.ru/" theme="dark">
        <div class="wrapper-dark">
          <div class="error-view">
            <div class="error-view__message">
              Unknown API error
            </div>
            <div class="error-view__details"></div>
          </div>
        </div>
        <div class="error-view"></div>
        </identifo-form>
      `);
  });
  // it('renders with values', async () => {
  //   const { root } = await newSpecPage({
  //     components: [IdentifoForm],
  //     html: `<identifo-form first="Stencil" last="'Don't call me a framework' JS"></identifo-form>`,
  //   });
  //   expect(root).toEqualHtml(`
  //     <identifo-form first="Stencil" last="'Don't call me a framework' JS">
  //       <mock:shadow-root>
  //         <div>
  //           Hello, World! I'm Stencil 'Don't call me a framework' JS
  //         </div>
  //       </mock:shadow-root>
  //     </identifo-form>
  //   `);
  // });
});
