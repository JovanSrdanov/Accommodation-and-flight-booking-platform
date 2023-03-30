import { useEffect } from 'react';

const useHoverAnimation = () => {
    const uppercaseLetters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const lowercaseLetters = 'abcdefghijklmnopqrstuvwxyz';
    useEffect(() => {
        const headers = document.querySelectorAll('h1');
        headers.forEach((header) => {
            let interval = null;
            const originalText = header.innerText.trim();
            header.dataset.originalText = originalText;
            header.addEventListener('mouseover', (event) => {
                const text = event.target.dataset.originalText;
                let iteration = 0;

                clearInterval(interval);

                interval = setInterval(() => {
                    event.target.innerText = text
                        .split('')
                        .map((letter, index) => {
                            if (letter === ' ') {
                                return letter;
                            }

                            if (index < iteration) {
                                return text[index];
                            }

                            if (uppercaseLetters.includes(letter)) {
                                return uppercaseLetters[Math.floor(Math.random() * 26)];
                            }

                            if (lowercaseLetters.includes(letter)) {
                                return lowercaseLetters[Math.floor(Math.random() * 26)];
                            }

                            return letter;
                        })
                        .join('');

                    if (iteration >= text.replace(/\s/g, '').length) {
                        clearInterval(interval);
                    }

                    iteration += 1 / 3;
                }, 30);
            });
        });
    }, []);
};

export default useHoverAnimation;
