import { FC } from 'react';
import Intro from './intro';


const Main: FC = () => {
    return (
        <main>
            <div className="container">
                <div className="main">
                    <Intro />
                </div>
            </div>
        </main>
    );
}

export default Main;