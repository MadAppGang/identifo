import { newSpecPage } from '@stencil/core/testing';
import { IdentifoForm } from './identifo-form';

describe('identifo-form', () => {
  it('renders', async () => {
    const { root } = await newSpecPage({
      components: [IdentifoForm],
      html: '<identifo-form></identifo-form>',
    });
    expect(root).toEqualHtml(`
      <identifo-form>
        <mock:shadow-root>
          <div>
            Hello, World! I'm
          </div>
        </mock:shadow-root>
      </identifo-form>
    `);
  });

  it('renders with values', async () => {
    const { root } = await newSpecPage({
      components: [IdentifoForm],
      html: `<identifo-form first="Stencil" last="'Don't call me a framework' JS"></identifo-form>`,
    });
    expect(root).toEqualHtml(`
      <identifo-form first="Stencil" last="'Don't call me a framework' JS">
        <mock:shadow-root>
          <div>
            Hello, World! I'm Stencil 'Don't call me a framework' JS
          </div>
        </mock:shadow-root>
      </identifo-form>
    `);
  });
});
