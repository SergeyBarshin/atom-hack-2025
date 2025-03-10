from flask import Flask, request, jsonify
import numpy as np
import requests
import math
from scipy.integrate import solve_ivp

app = Flask(__name__)


def get_pressures(n):
    pressures = []
    for i in range(n):
        url = f"https://yand.dyndns.org/api/hackathon.aspx?id={i+1}"
        response = requests.get(url)
        if response.status_code == 200:
            pressures.append(response.json().get("pressure", 0))
        else:
            pressures.append(0)  # Обработка ошибки запроса
    return np.array(pressures)


def calculate_liquid_level(volumes, pressures, liquid_level_init, matrix):
    # Параметры
    rho = 1000.0  # плотность воды, кг/м³
    g = 9.81  # ускорение свободного падения, м/с²

    # Расчет эффективного напора
    N = len(liquid_level_init)  # количество резервуаров
    h = np.array(liquid_level_init, dtype=float)  # начальные уровни жидкости
    p = np.array(pressures, dtype=float)  # давления
    A_res = (
        np.array(volumes, dtype=float) / h
    )  # площади сечения резервуаров (A = V / h)

    # Матрица соединений (A_conn)
    A_conn = np.array(matrix, dtype=float) * 0.01  # преобразуем площадь из дм² в м²

    # Эффективный напор
    H = h + p / (rho * g)  # эффективный напор H = h + p / (rho * g)

    # Потоки между резервуарами
    F = np.zeros(N, dtype=float)  # Массив для хранения потоков

    for i in range(N):
        for j in range(N):  # Проверяем все соединения
            if i != j and A_conn[i, j] > 0:
                head_diff = H[i] - H[j]  # Разница напоров
                if abs(head_diff) < 1e-6:  # Если разница напоров мала, поток нулевой
                    Q = 0.0
                else:
                    Q = A_conn[i, j] * math.sqrt(
                        2 * g * abs(head_diff)
                    )  # Расчет потока

                # Направление потока
                if head_diff > 0:
                    F[i] -= Q
                    F[j] += Q
                elif head_diff < 0:
                    F[i] += Q
                    F[j] -= Q

    # Изменение уровней жидкости
    dh_dt = F / A_res  # Изменение уровня жидкости (dh/dt = net_flow / A_res)
    h_new = h + dh_dt  # Новый уровень жидкости

    # Обрезаем уровень жидкости, чтобы он не выходил за пределы объёмов резервуаров
    h_new = np.clip(h_new, 0, volumes)  # Уровни не могут превышать объём резервуара

    return h_new


@app.route("/update", methods=["POST"])
def process():
    data = request.get_json()
    volumes = np.array(data.get("volumes", []))
    matrix = np.array(data.get("matrix", []))
    liquid_level_init = np.array(data.get("liquid_level", []))

    if volumes.size == 0 or matrix.size == 0 or liquid_level_init.size == 0:
        return jsonify({"error": "Invalid input data"}), 400

    n = volumes.size

    pressures = get_pressures(n)

    h_new = calculate_liquid_level(volumes, pressures, liquid_level_init, matrix)

    data["liquid_level"] = h_new.tolist()
    data["pressures"] = list(pressures)

    return jsonify(data)


if __name__ == "__main__":
    app.run(debug=True)
