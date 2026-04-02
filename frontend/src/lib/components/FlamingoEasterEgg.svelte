<script lang="ts">
    import { onDestroy } from "svelte";

    type Flamingo = {
        id: number;
        isGolden: boolean;
        left: number;
        size: number;
        peak: number;
        fall: number;
        duration: number;
        driftStart: number;
        driftMid: number;
        driftEnd: number;
        driftOut: number;
        rotateStart: number;
        rotateMid: number;
        rotateEnd: number;
        rotateOut: number;
    };

    let { trigger = 0 }: { trigger: number } = $props();

    let flock: Flamingo[] = $state([]);
    let lastTrigger = 0;
    let nextId = 0;
    const cleanupTimers = new Map<number, number>();

    function randomBetween(min: number, max: number): number {
        return min + Math.random() * (max - min);
    }

    function spawnFlamingo() {
        const id = nextId++;
        const isGolden = Math.random() < 0.01;
        const particle: Flamingo = {
            id,
            isGolden,
            left: randomBetween(10, 90),
            size: randomBetween(isGolden ? 30 : 26, isGolden ? 48 : 42),
            peak: randomBetween(180, 340),
            fall: randomBetween(90, 150),
            duration: randomBetween(1800, 2600),
            driftStart: randomBetween(-18, 18),
            driftMid: randomBetween(-72, 72),
            driftEnd: randomBetween(-48, 48),
            driftOut: randomBetween(-96, 96),
            rotateStart: randomBetween(-22, 14),
            rotateMid: randomBetween(-8, 26),
            rotateEnd: randomBetween(-16, 18),
            rotateOut: randomBetween(-28, 28),
        };

        flock = [...flock, particle];

        const timer = window.setTimeout(() => {
            flock = flock.filter((entry) => entry.id !== id);
            cleanupTimers.delete(id);
        }, particle.duration);

        cleanupTimers.set(id, timer);
    }

    $effect(() => {
        if (trigger <= lastTrigger) return;

        const launches = Math.min(trigger - lastTrigger, 4);
        lastTrigger = trigger;

        for (let i = 0; i < launches; i++) {
            spawnFlamingo();
        }
    });

    onDestroy(() => {
        for (const timer of cleanupTimers.values()) {
            window.clearTimeout(timer);
        }
        cleanupTimers.clear();
    });
</script>

<div class="flamingo-overlay" aria-hidden="true">
    {#each flock as flamingo (flamingo.id)}
        <span
            class="flamingo-shell"
            class:flamingo-shell--golden={flamingo.isGolden}
            style={`left:${flamingo.left}%; font-size:${flamingo.size}px; --peak:${flamingo.peak}px; --fall:${flamingo.fall}px; --duration:${flamingo.duration}ms; --drift-start:${flamingo.driftStart}px; --drift-mid:${flamingo.driftMid}px; --drift-end:${flamingo.driftEnd}px; --drift-out:${flamingo.driftOut}px; --rotate-start:${flamingo.rotateStart}deg; --rotate-mid:${flamingo.rotateMid}deg; --rotate-end:${flamingo.rotateEnd}deg; --rotate-out:${flamingo.rotateOut}deg;`}
        >
            {#if flamingo.isGolden}
                <span class="flamingo-aura"></span>
                <span class="flamingo-sparkle flamingo-sparkle--one">✦</span>
                <span class="flamingo-sparkle flamingo-sparkle--two">✦</span>
                <span class="flamingo-sparkle flamingo-sparkle--three">✦</span>
            {/if}
            <span class="flamingo" class:flamingo--golden={flamingo.isGolden}>
                🦩
            </span>
        </span>
    {/each}
</div>

<style>
    .flamingo-overlay {
        position: absolute;
        inset: 0;
        overflow: hidden;
        pointer-events: none;
        z-index: 40;
    }

    .flamingo-shell {
        position: absolute;
        bottom: 10px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        line-height: 1;
        transform: translate3d(0, 120%, 0);
        transform-origin: center;
        animation: flamingo-flight var(--duration) cubic-bezier(0.18, 0.8, 0.22, 1)
            forwards;
        will-change: transform, opacity;
        opacity: 0;
    }

    .flamingo {
        position: relative;
        z-index: 2;
        line-height: 1;
        filter: drop-shadow(0 8px 10px rgba(0, 0, 0, 0.32));
    }

    .flamingo--golden {
        filter:
            saturate(1.25)
            hue-rotate(-24deg)
            brightness(1.12)
            drop-shadow(0 0 8px rgba(255, 208, 92, 0.55))
            drop-shadow(0 0 18px rgba(255, 186, 48, 0.32))
            drop-shadow(0 10px 14px rgba(0, 0, 0, 0.3));
        animation: flamingo-golden-pulse 640ms ease-in-out infinite alternate;
    }

    .flamingo-aura {
        position: absolute;
        inset: 50% auto auto 50%;
        width: 1.9em;
        height: 1.9em;
        transform: translate(-50%, -50%);
        border-radius: 999px;
        background:
            radial-gradient(
                circle,
                rgba(255, 221, 128, 0.55) 0%,
                rgba(255, 202, 64, 0.26) 40%,
                rgba(255, 202, 64, 0) 72%
            );
        filter: blur(4px);
        animation: flamingo-aura-pulse 700ms ease-in-out infinite alternate;
    }

    .flamingo-sparkle {
        position: absolute;
        z-index: 1;
        font-size: 0.38em;
        line-height: 1;
        color: #ffe8a3;
        text-shadow:
            0 0 6px rgba(255, 232, 163, 0.9),
            0 0 12px rgba(255, 193, 79, 0.6);
        animation: flamingo-sparkle 880ms ease-in-out infinite;
    }

    .flamingo-sparkle--one {
        top: -0.15em;
        left: -0.35em;
        animation-delay: 0ms;
    }

    .flamingo-sparkle--two {
        top: 0.1em;
        right: -0.42em;
        animation-delay: 180ms;
    }

    .flamingo-sparkle--three {
        bottom: 0;
        right: 0.1em;
        animation-delay: 360ms;
    }

    @keyframes flamingo-flight {
        0% {
            opacity: 0;
            transform: translate3d(var(--drift-start), 120%, 0)
                rotate(var(--rotate-start)) scale(0.92);
        }

        12% {
            opacity: 1;
        }

        34% {
            opacity: 1;
            transform: translate3d(var(--drift-mid), calc(var(--peak) * -1), 0)
                rotate(var(--rotate-mid)) scale(1.08);
        }

        68% {
            opacity: 1;
            transform: translate3d(var(--drift-end), calc(var(--fall) * -1), 0)
                rotate(var(--rotate-end)) scale(0.98);
        }

        100% {
            opacity: 0;
            transform: translate3d(var(--drift-out), 135%, 0)
                rotate(var(--rotate-out)) scale(0.88);
        }
    }

    @keyframes flamingo-golden-pulse {
        from {
            transform: scale(1);
        }

        to {
            transform: scale(1.08);
        }
    }

    @keyframes flamingo-aura-pulse {
        from {
            opacity: 0.62;
            transform: translate(-50%, -50%) scale(0.92);
        }

        to {
            opacity: 0.95;
            transform: translate(-50%, -50%) scale(1.08);
        }
    }

    @keyframes flamingo-sparkle {
        0%,
        100% {
            opacity: 0.15;
            transform: scale(0.7) translateY(1px);
        }

        50% {
            opacity: 1;
            transform: scale(1.2) translateY(-2px);
        }
    }

    @media (prefers-reduced-motion: reduce) {
        .flamingo-shell,
        .flamingo,
        .flamingo-aura,
        .flamingo-sparkle {
            animation-duration: 1ms;
        }
    }
</style>
