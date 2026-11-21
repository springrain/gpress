(function () {
  "use strict";

  var host = document.querySelector("[data-footer-water]");
  if (!host) {
    return;
  }

  var canvas = host.querySelector("canvas");
  var context = canvas && canvas.getContext("2d");
  if (!context) {
    return;
  }

  var footer = host.closest(".tm-footer") || host;
  var root = document.documentElement;
  var motionQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
  var width = 0;
  var height = 0;
  var surfaceY = 0;
  var points = [];
  var leftDeltas = [];
  var rightDeltas = [];
  var fishes = [];
  var droplets = [];
  var colors = {};
  var frameId = 0;
  var lastTime = 0;
  var accumulator = 0;
  var spawnDelay = 0;
  var schoolDelay = 0;
  var visible = true;
  var disposed = false;
  var STEP = 1 / 60;
  var fishPaths = createFishPaths();

  function clamp(value, min, max) {
    return Math.max(min, Math.min(max, value));
  }

  function random(min, max) {
    return min + Math.random() * (max - min);
  }

  function readColors() {
    var style = window.getComputedStyle(footer);
    colors.water = style.getPropertyValue("--tm-footer-water-color").trim() || "rgba(72, 146, 194, 0.14)";
    colors.fish = style.getPropertyValue("--tm-footer-fish-color").trim() || "rgba(35, 118, 183, 0.34)";
  }

  function createFishPaths() {
    if (typeof Path2D !== "function") {
      return null;
    }

    var tail = new Path2D();
    tail.moveTo(0, 0);
    tail.bezierCurveTo(-2.2, -1, -4.2, -4.8, -7.4, -6.2);
    tail.bezierCurveTo(-6.2, -2.6, -4.1, -0.8, -2, 0);
    tail.bezierCurveTo(-4.2, 0.8, -6.2, 2.6, -7.4, 6.2);
    tail.bezierCurveTo(-4.2, 4.8, -2.2, 1, 0, 0);
    tail.closePath();

    var fins = new Path2D();
    fins.moveTo(4, -2.6);
    fins.bezierCurveTo(1.8, -4.8, -0.7, -6.1, -3.2, -5.2);
    fins.bezierCurveTo(-1.2, -3.4, 0.8, -2.4, 3.2, -1.8);
    fins.closePath();
    fins.moveTo(3.2, 2.5);
    fins.bezierCurveTo(1.2, 5.2, -1.5, 7, -4.2, 6.2);
    fins.bezierCurveTo(-2.1, 3.8, 0.2, 2.6, 2.5, 1.8);
    fins.closePath();

    return { tail: tail, fins: fins };
  }

  function createSurface() {
    var spacing = width < 680 ? 10 : 12;
    var count = Math.max(3, Math.ceil(width / spacing) + 1);
    points = [];

    for (var index = 0; index < count; index += 1) {
      points.push({
        x: (index / (count - 1)) * width,
        y: surfaceY,
        velocity: 0
      });
    }
    leftDeltas = new Array(count).fill(0);
    rightDeltas = new Array(count).fill(0);
  }

  function resize() {
    var bounds = host.getBoundingClientRect();
    var nextWidth = Math.max(1, Math.round(bounds.width));
    var nextHeight = Math.max(1, Math.round(bounds.height));
    if (nextWidth === width && nextHeight === height) {
      return;
    }

    width = nextWidth;
    height = nextHeight;
    surfaceY = Math.round(height * 0.52);

    var pixelRatio = clamp(window.devicePixelRatio || 1, 1, 2);
    canvas.width = Math.round(width * pixelRatio);
    canvas.height = Math.round(height * pixelRatio);
    context.setTransform(pixelRatio, 0, 0, pixelRatio, 0, 0);

    createSurface();
    fishes = [];
    droplets = [];
    spawnDelay = 0;
    schoolDelay = random(7, 11);
    render();
  }

  function disturbSurface(x, impulse) {
    if (!points.length || x < 0 || x > width) {
      return;
    }

    var center = Math.round((x / width) * (points.length - 1));
    for (var offset = -2; offset <= 2; offset += 1) {
      var point = points[center + offset];
      if (point) {
        point.velocity += impulse * (1 - Math.abs(offset) / 3);
      }
    }
  }

  function updateSurface() {
    points.forEach(function (point) {
      point.velocity += (surfaceY - point.y) * 0.024;
      point.velocity *= 0.94;
      point.y += point.velocity;
    });

    for (var pass = 0; pass < 2; pass += 1) {
      for (var index = 0; index < points.length; index += 1) {
        if (index > 0) {
          leftDeltas[index] = 0.105 * (points[index].y - points[index - 1].y);
          points[index - 1].velocity += leftDeltas[index];
        }
        if (index < points.length - 1) {
          rightDeltas[index] = 0.105 * (points[index].y - points[index + 1].y);
          points[index + 1].velocity += rightDeltas[index];
        }
      }

      for (var applyIndex = 0; applyIndex < points.length; applyIndex += 1) {
        if (applyIndex > 0) {
          points[applyIndex - 1].y += leftDeltas[applyIndex];
        }
        if (applyIndex < points.length - 1) {
          points[applyIndex + 1].y += rightDeltas[applyIndex];
        }
      }
    }
  }

  function createFish(delay) {
    var direction = Math.random() < 0.5 ? 1 : -1;
    var speed = random(125, 210);
    return {
      x: direction > 0 ? -42 : width + 42,
      y: surfaceY + random(20, Math.max(28, height - surfaceY - 16)),
      vx: direction * speed,
      vy: random(-10, 10),
      direction: direction,
      scale: random(0.85, 1.08),
      phase: random(0, Math.PI * 2),
      phaseSpeed: random(5.4, 7.6),
      swimDepth: random(18, Math.max(24, height - surfaceY - 20)),
      buoyancy: random(175, 235),
      jumpVelocity: random(205, 275),
      jumpCooldown: random(0.2, 3.8),
      jumping: false,
      fullyAbove: false,
      angle: 0,
      age: -delay,
    };
  }

  function desiredFishCount() {
    if (width < 520) {
      return 1;
    }
    return clamp(Math.round(width / 430), 2, 4);
  }

  function emitDroplets(fish, tailX) {
    var count = Math.random() < 0.7 ? 2 : 3;
    for (var index = 0; index < count; index += 1) {
      var life = random(0.42, 0.68);
      droplets.push({
        x: tailX + random(-2, 2),
        y: surfaceY + random(-1, 1),
        vx: fish.vx * 0.08 + random(-22, 22),
        vy: random(-76, -38),
        radius: random(0.55, 1.05),
        life: life,
        maxLife: life
      });
    }
  }

  function updateDroplets(delta) {
    droplets = droplets.filter(function (drop) {
      drop.life -= delta;
      drop.x += drop.vx * delta;
      drop.y += drop.vy * delta;
      drop.vy += 190 * delta;
      return drop.life > 0 && drop.y < height + 4;
    });
  }

  function updateFish(fish, delta) {
    fish.age += delta;
    if (fish.age < 0) {
      return true;
    }

    fish.x += fish.vx * delta;
    fish.phase += delta * fish.phaseSpeed;

    if (fish.y < surfaceY) {
      fish.vy += 420 * delta;
    } else if (fish.jumping) {
      fish.vy -= fish.buoyancy * delta;
    } else {
      fish.jumpCooldown -= delta;
      var targetY = surfaceY + fish.swimDepth;
      fish.vy += (targetY - fish.y) * 3.2 * delta;
      fish.vy *= Math.pow(0.92, delta * 60);
      if (fish.jumpCooldown <= 0) {
        fish.jumping = true;
        fish.vy = -fish.jumpVelocity;
      }
    }
    fish.vy = clamp(fish.vy, -305, 300);
    fish.y += fish.vy * delta;

    var desiredAngle = Math.atan2(fish.vy, Math.abs(fish.vx)) * fish.direction;
    fish.angle += (desiredAngle - fish.angle) * 0.16;

    var headY = fish.y + Math.sin(fish.angle) * 13 * fish.scale * fish.direction;
    var tailY = fish.y - Math.sin(fish.angle) * 18 * fish.scale * fish.direction;
    var tailX = fish.x - Math.cos(fish.angle) * 18 * fish.scale * fish.direction;
    var fullyAbove = Math.max(headY, tailY) < surfaceY;
    var fullyBelow = Math.min(headY, tailY) >= surfaceY;

    if (!fish.fullyAbove && fullyAbove) {
      fish.fullyAbove = true;
      disturbSurface(tailX, clamp(fish.vy * 0.024, -7, 7));
      emitDroplets(fish, tailX);
    } else if (fish.fullyAbove && fullyBelow) {
      fish.fullyAbove = false;
      fish.jumping = false;
      fish.jumpCooldown = random(1.4, 5.2);
      disturbSurface(tailX, clamp(fish.vy * 0.055, 9, 16));
    }

    return fish.direction > 0 ? fish.x < width + 48 : fish.x > -48;
  }

  function scheduleSchoolJump(delta) {
    schoolDelay -= delta;
    if (schoolDelay > 0 || fishes.length < 2) {
      return;
    }

    var candidates = fishes.filter(function (fish) {
      return fish.age >= 0 && !fish.jumping && fish.y >= surfaceY;
    });
    var count = Math.min(candidates.length, Math.random() < 0.65 ? 2 : 3);
    for (var index = 0; index < count; index += 1) {
      var selected = candidates.splice(Math.floor(Math.random() * candidates.length), 1)[0];
      selected.jumpCooldown = random(0, 0.3);
    }
    schoolDelay = random(7, 13);
  }

  function update(delta) {
    updateSurface();
    updateDroplets(delta);
    scheduleSchoolJump(delta);

    fishes = fishes.filter(function (fish) {
      return updateFish(fish, delta);
    });

    spawnDelay -= delta;
    if (fishes.length < desiredFishCount() && spawnDelay <= 0) {
      fishes.push(createFish(fishes.length * 0.45));
      spawnDelay = random(0.65, 1.35);
    }
  }

  function drawWater() {
    context.save();
    context.fillStyle = colors.water;
    context.beginPath();
    context.moveTo(0, height);
    context.lineTo(points[0].x, points[0].y);

    for (var index = 1; index < points.length - 1; index += 1) {
      var point = points[index];
      var next = points[index + 1];
      context.quadraticCurveTo(point.x, point.y, (point.x + next.x) / 2, (point.y + next.y) / 2);
    }

    var last = points[points.length - 1];
    context.lineTo(last.x, last.y);
    context.lineTo(width, height);
    context.closePath();
    context.fill();
    context.restore();
  }

  function drawTail() {
    if (fishPaths) {
      context.fill(fishPaths.tail);
      return;
    }
    context.beginPath();
    context.moveTo(0, 0);
    context.bezierCurveTo(-2.2, -1, -4.2, -4.8, -7.4, -6.2);
    context.bezierCurveTo(-6.2, -2.6, -4.1, -0.8, -2, 0);
    context.bezierCurveTo(-4.2, 0.8, -6.2, 2.6, -7.4, 6.2);
    context.bezierCurveTo(-4.2, 4.8, -2.2, 1, 0, 0);
    context.closePath();
    context.fill();
  }

  function drawFins() {
    if (fishPaths) {
      context.fill(fishPaths.fins);
      return;
    }
    context.beginPath();
    context.moveTo(4, -2.6);
    context.bezierCurveTo(1.8, -4.8, -0.7, -6.1, -3.2, -5.2);
    context.bezierCurveTo(-1.2, -3.4, 0.8, -2.4, 3.2, -1.8);
    context.closePath();
    context.moveTo(3.2, 2.5);
    context.bezierCurveTo(1.2, 5.2, -1.5, 7, -4.2, 6.2);
    context.bezierCurveTo(-2.1, 3.8, 0.2, 2.6, 2.5, 1.8);
    context.closePath();
    context.fill();
  }

  function drawFish(fish) {
    if (fish.age < 0) {
      return;
    }

    var speedRatio = clamp((Math.abs(fish.vx) - 125) / 85, 0, 1);
    var tailWave = Math.sin(fish.phase);
    var airFactor = fish.y < surfaceY ? 0.65 : 1;
    var tailAngle = tailWave * (12 + speedRatio * 6) * airFactor;
    var bodyBend = tailWave * (0.65 + speedRatio * 0.95);
    if (fish.y < surfaceY) {
      bodyBend += clamp(fish.vy / 280, -1, 1) * 0.8;
    }

    context.save();
    context.translate(fish.x, fish.y);
    context.rotate(fish.angle);
    context.scale(fish.direction * fish.scale, fish.scale);
    context.fillStyle = colors.fish;

    // The stable head leads a softly flexing rear body.
    context.beginPath();
    context.moveTo(-12, bodyBend);
    context.bezierCurveTo(-8.5, -2.8 + bodyBend * 0.65, -2.3, -4.6, 4, -4.4);
    context.bezierCurveTo(9.5, -4.2, 13.2, -2.4, 14, 0);
    context.bezierCurveTo(13.2, 2.5, 9.5, 4.3, 4, 4.5);
    context.bezierCurveTo(-2.3, 4.7, -8.5, 3 + bodyBend * 0.65, -12, bodyBend);
    context.closePath();
    context.fill();

    drawFins();

    context.save();
    context.translate(-11.5, bodyBend);
    context.rotate(tailAngle * Math.PI / 180);
    drawTail();
    context.restore();
    context.restore();
  }

  function drawDroplets() {
    context.save();
    context.fillStyle = colors.fish;
    droplets.forEach(function (drop) {
      context.globalAlpha = clamp(drop.life / drop.maxLife, 0, 1) * 0.62;
      context.beginPath();
      context.arc(drop.x, drop.y, drop.radius, 0, Math.PI * 2);
      context.fill();
    });
    context.restore();
  }

  function render() {
    if (!width || !height || !points.length) {
      return;
    }

    context.clearRect(0, 0, width, height);
    fishes.forEach(drawFish);
    drawDroplets();
    drawWater();
  }

  function animate(time) {
    frameId = 0;
    if (disposed || !visible || document.hidden || motionQuery.matches) {
      return;
    }

    if (!lastTime) {
      lastTime = time;
    }
    var elapsed = Math.min((time - lastTime) / 1000, 0.05);
    lastTime = time;
    accumulator += elapsed;

    while (accumulator >= STEP) {
      update(STEP);
      accumulator -= STEP;
    }

    render();
    frameId = window.requestAnimationFrame(animate);
  }

  function start() {
    if (disposed || frameId || !visible || document.hidden || motionQuery.matches) {
      render();
      return;
    }
    lastTime = 0;
    accumulator = 0;
    frameId = window.requestAnimationFrame(animate);
  }

  function stop() {
    if (frameId) {
      window.cancelAnimationFrame(frameId);
      frameId = 0;
    }
  }

  function handleMotionChange() {
    if (motionQuery.matches) {
      stop();
      fishes = [];
      droplets = [];
      createSurface();
      render();
    } else {
      start();
    }
  }

  readColors();
  resize();

  var resizeObserver = typeof ResizeObserver === "function" ? new ResizeObserver(resize) : null;
  if (resizeObserver) {
    resizeObserver.observe(host);
  } else {
    window.addEventListener("resize", resize);
  }

  var intersectionObserver = typeof IntersectionObserver === "function" ? new IntersectionObserver(function (entries) {
    visible = entries[0].isIntersecting;
    if (visible) {
      start();
    } else {
      stop();
    }
  }, { rootMargin: "160px 0px" }) : null;
  if (intersectionObserver) {
    intersectionObserver.observe(footer);
  }

  var themeObserver = typeof MutationObserver === "function" ? new MutationObserver(function () {
    readColors();
    render();
  }) : null;
  if (themeObserver) {
    themeObserver.observe(root, { attributes: true, attributeFilter: ["class"] });
  }

  document.addEventListener("visibilitychange", function () {
    if (document.hidden) {
      stop();
    } else {
      start();
    }
  });

  if (typeof motionQuery.addEventListener === "function") {
    motionQuery.addEventListener("change", handleMotionChange);
  } else if (typeof motionQuery.addListener === "function") {
    motionQuery.addListener(handleMotionChange);
  }

  window.addEventListener("pagehide", function () {
    disposed = true;
    stop();
    if (resizeObserver) {
      resizeObserver.disconnect();
    }
    if (intersectionObserver) {
      intersectionObserver.disconnect();
    }
    if (themeObserver) {
      themeObserver.disconnect();
    }
  }, { once: true });

  handleMotionChange();
})();
