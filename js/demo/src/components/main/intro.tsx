import React, { FC } from 'react';

type IntroProps = {

}

const Intro: FC<IntroProps> = () => {
    return (
        <div className="animated-title">
            <div className="text-top">
                <div>
                    <span>Welcome</span>
                    <span>to the identifo.js</span>
                </div>
            </div>
            <div className="text-bottom">
                <div>@nikita-ks</div>
            </div>
        </div>
    );
}

export default Intro;