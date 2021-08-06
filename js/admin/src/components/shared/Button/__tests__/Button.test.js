import Button from '..';

describe('<Button />', () => {
  test('renders as expected', () => {
    expect(shallow(<Button />)).toMatchSnapshot();
  });

  test('renders children when passed in', () => {
    const output = shallow(<Button><div>Hello, world!</div></Button>);
    expect(output.contains(<div>Hello, world!</div>)).toBe(true);
  });

  test('calls onClick prop function on click', () => {
    const onClick = jest.fn();

    shallow(<Button onClick={onClick}>Click me</Button>)
      .find('button')
      .simulate('click');

    expect(onClick).toHaveBeenCalled();
  });

  test('renders as expected when disabled', () => {
    expect(shallow(<Button disabled>Text</Button>)).toMatchSnapshot();
  });
});
