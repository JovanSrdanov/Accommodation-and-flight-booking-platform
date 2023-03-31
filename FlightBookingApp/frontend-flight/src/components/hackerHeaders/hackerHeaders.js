import {useEffect} from "react";

function HackerHeaders() {
    useEffect(() => {
        function handleHover(event) {
            const uppercaseLetters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
            const lowercaseLetters = 'abcdefghijklmnopqrstuvwxyz';

            const {target} = event;
            if (target.tagName === 'H1') {
                let interval = null;
                const originalText = target.dataset.originalText || target.innerText.trim();

                const text = originalText;
                let iteration = 0;
                let animationComplete = false;
                let animationStarted = false;
                let animationReversed = false;

                const originalTextWithoutSpaces = originalText.replace(/\s/g, '');

                clearInterval(interval);

                interval = setInterval(() => {
                    let newText = '';
                    for (let i = 0; i < text.length; i++) {
                        const letter = text[i];
                        if (letter === ' ') {
                            newText += letter;
                            continue;
                        }

                        if (i < iteration) {
                            newText += text[i];
                        } else {
                            if (!animationStarted) {
                                animationStarted = true;
                                target.dataset.originalText = originalText;
                            }

                            let newLetter;
                            if (uppercaseLetters.includes(letter)) {
                                newLetter = uppercaseLetters[Math.floor(Math.random() * 26)];
                            } else if (lowercaseLetters.includes(letter)) {
                                newLetter = lowercaseLetters[Math.floor(Math.random() * 26)];
                            } else {
                                newLetter = letter;
                            }

                            if (animationComplete) {
                                if (newLetter !== originalText[i]) {
                                    animationReversed = true;
                                }
                                newLetter = originalText[i];
                            }
                            newText += newLetter;
                        }
                    }

                    target.innerText = newText;

                    if (iteration >= originalTextWithoutSpaces.length) {
                        clearInterval(interval);
                        animationComplete = true;
                        target.innerText = originalText;
                    }

                    if (animationComplete && animationStarted && !animationReversed) {
                        iteration -= 0.5;
                    } else {
                        iteration += 0.5;
                    }
                }, 50);
            }
        }

        document.addEventListener('mouseover', handleHover);
        return () => {
            document.removeEventListener('mouseover', handleHover);
        };
    }, []);
}

export default HackerHeaders